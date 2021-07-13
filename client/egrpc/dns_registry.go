package egrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
)

// NOTICE: 不支持ServiceConfig
type dnsRegistry struct {
	r                    dnsResolver
	nameWrapper          func(string) string
	forceResolveInterval time.Duration
}

// DNSRegistryOption ...
type DNSRegistryOption func(c *dnsRegistry)

// WithNameWrapper 对dns:///SVC-NAME:PORT中的SVC-NAME进行处理，当SVC-NAME为IP时不做任何处理
func WithNameWrapper(wrapper func(string) string) DNSRegistryOption {
	return func(c *dnsRegistry) {
		c.nameWrapper = wrapper
	}
}

// WithForceResolveInterval 指定DNS强制resolve周期
func WithForceResolveInterval(forceResolveInterval time.Duration) DNSRegistryOption {
	return func(c *dnsRegistry) {
		c.forceResolveInterval = forceResolveInterval
	}
}

// DNSRegistry 返回DNS Registry
func DNSRegistry(opts ...DNSRegistryOption) eregistry.Registry {
	r := &dnsRegistry{
		nameWrapper:          func(name string) string { return name },
		forceResolveInterval: minDNSResRate,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (d *dnsRegistry) resolve(ctx context.Context, host, port, schema string) (*eregistry.Endpoints, error) {
	state, err := d.r.lookup(ctx, host, port)
	if err != nil {
		return nil, err
	}
	endpoints := &eregistry.Endpoints{
		Nodes: make(map[string]server.ServiceInfo),
	}
	for _, addr := range state.Addresses {
		endpoints.Nodes[addr.Addr] = server.ServiceInfo{Address: addr.Addr}
	}
	return endpoints, nil
}

func (d *dnsRegistry) WatchServices(ctx context.Context, target eregistry.Target) (chan eregistry.Endpoints, error) {
	if target.Authority == "" {
		d.r = dnsResolver{
			resolver: defaultResolver,
			rn:       make(chan struct{}, 1),
		}
	} else {
		ar, err := customAuthorityResolver(target.Authority)
		if err != nil {
			return nil, err
		}
		d.r = dnsResolver{
			resolver: ar,
			rn:       make(chan struct{}, 1),
		}
	}

	// 从endpoint中解析出host和port
	host, port, err := parseTarget(target.Endpoint, defaultPort)
	if err != nil {
		return nil, err
	}
	var endpointsCh = make(chan eregistry.Endpoints, 10)

	// 如果是静态IP，则直接返回静态IP
	if ipAddr, ok := formatIP(host); ok {
		addr := ipAddr + ":" + port
		endpointsCh <- eregistry.Endpoints{Nodes: map[string]server.ServiceInfo{addr: {Address: addr}}}
		return endpointsCh, nil
	}

	// 如果是域名，则尝试解析域名
	host = d.nameWrapper(host)
	endpoints, err := d.resolve(ctx, host, port, target.Scheme)
	if err != nil {
		return nil, err
	}
	endpointsCh <- *endpoints.DeepCopy()

	// 开启定时器，定时重新解析域名，解决扩容时无法发现新实例问题
	go func() {
		ticker := time.NewTicker(d.forceResolveInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-d.r.rn:
				endpointsCh <- *d.mustResolve(ctx, host, port, target).DeepCopy()
			case <-ticker.C:
				endpointsCh <- *d.mustResolve(ctx, host, port, target).DeepCopy()
			}
		}
	}()
	return endpointsCh, err
}

func (d *dnsRegistry) mustResolve(ctx context.Context, host, port string, target eregistry.Target) *eregistry.Endpoints {
	newEndpoints, err := d.resolve(ctx, host, port, target.Scheme)
	if err != nil {
		elog.Warn("resolve failed", elog.FieldErr(err))
	}
	return newEndpoints
}

func (d *dnsRegistry) ListServices(ctx context.Context, target eregistry.Target) (services []*server.ServiceInfo, err error) {
	return nil, nil
}

// RegisterService do noting
func (d *dnsRegistry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	panic("not implement RegisterService yet")
}

// UnregisterService do noting
func (d *dnsRegistry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	panic("not implement UnregisterService yet")
}

func (d *dnsRegistry) Close() error {
	return nil
}

func (d *dnsRegistry) SyncServices(ctx context.Context, opts eregistry.SyncServicesOptions) error {
	d.r.ResolveNow(opts.GrpcResolverNowOptions)
	return nil
}

// EnableSRVLookups controls whether the DNS resolver attempts to fetch gRPCLB
// copy from google.golang.org/grpc@v1.29.1/resolver/dns/dns_resolver.go
// addresses from SRV records.  Must not be changed after init time.
var EnableSRVLookups = false

const (
	defaultPort       = "443"
	defaultDNSSvrPort = "53"
)

var (
	errMissingAddr = errors.New("dns resolver: missing address")

	// Addresses ending with a colon that is supposed to be the separator
	// between host and port is not allowed.  E.g. "::" is a valid address as
	// it is an IPv6 address (host only) and "[::]:" is invalid as it ends with
	// a colon as the host and port separator
	errEndsWithColon = errors.New("dns resolver: missing port after port-separator colon")
)

var (
	defaultResolver netResolver = net.DefaultResolver
	// To prevent excessive re-resolution, we enforce a rate limit on DNS
	// resolution requests
	minDNSResRate = 120 * time.Second
)

type netResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
}

// dnsResolver watches for the name resolution update for a non-IP target.
type dnsResolver struct {
	resolver netResolver
	rn       chan struct{}
}

// ResolveNow invoke an immediate resolution of the target that this dnsResolver watches.
func (d *dnsResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}

func (d *dnsResolver) lookupSRV(ctx context.Context, host, port string) ([]resolver.Address, error) {
	if !EnableSRVLookups {
		return nil, nil
	}
	var newAddrs []resolver.Address
	_, srvs, err := d.resolver.LookupSRV(ctx, "grpclb", "tcp", host)
	if err != nil {
		err = handleDNSError(err, "SRV") // may become nil
		return nil, err
	}
	for _, s := range srvs {
		lbAddrs, err := d.resolver.LookupHost(ctx, s.Target)
		if err != nil {
			err = handleDNSError(err, "A") // may become nil
			if err == nil {
				// If there are other SRV records, look them up and ignore this
				// one that does not exist.
				continue
			}
			return nil, err
		}
		for _, a := range lbAddrs {
			ip, ok := formatIP(a)
			if !ok {
				return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
			}
			addr := ip + ":" + strconv.Itoa(int(s.Port))
			newAddrs = append(newAddrs, resolver.Address{Addr: addr, Type: resolver.GRPCLB, ServerName: s.Target}) //nolint
		}
	}
	return newAddrs, nil
}

var filterError = func(err error) error {
	if dnsErr, ok := err.(*net.DNSError); ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary {
		// Timeouts and temporary errors should be communicated to gRPC to
		// attempt another DNS query (with backoff).  Other errors should be
		// suppressed (they may represent the absence of a TXT record).
		return nil
	}
	return err
}

func handleDNSError(err error, lookupType string) error {
	err = filterError(err)
	if err != nil {
		err = fmt.Errorf("dns: %v record lookup error: %v", lookupType, err)
		grpclog.Infoln(err)
	}
	return err
}

func (d *dnsResolver) lookupHost(ctx context.Context, host, port string) ([]resolver.Address, error) {
	var newAddrs []resolver.Address
	addrs, err := d.resolver.LookupHost(ctx, host)
	if err != nil {
		err = handleDNSError(err, "A")
		return nil, err
	}
	for _, a := range addrs {
		ip, ok := formatIP(a)
		if !ok {
			return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
		}
		addr := ip + ":" + port
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}
	return newAddrs, nil
}

func (d *dnsResolver) lookup(ctx context.Context, host, port string) (*resolver.State, error) {
	srv, srvErr := d.lookupSRV(ctx, host, port)
	addrs, hostErr := d.lookupHost(ctx, host, port)
	if hostErr != nil && (srvErr != nil || len(srv) == 0) {
		return nil, hostErr
	}
	state := &resolver.State{
		Addresses: append(addrs, srv...),
	}
	// 不支持获取ServiceConfig
	return state, nil
}

// formatIP returns ok = false if addr is not a valid textual representation of an IP address.
// If addr is an IPv4 address, return the addr and ok = true.
// If addr is an IPv6 address, return the addr enclosed in square brackets and ok = true.
func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}

// parseTarget takes the user input target string and default port, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in square brackets, brackets
// are stripped when setting the host.
// examples:
// target: "www.google.com" defaultPort: "443" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" defaultPort: "443" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" defaultPort: "443" returns host: "ipv6-host", port: "443"
// target: ":80" defaultPort: "443" returns host: "localhost", port: "80"
func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", errMissingAddr
	}
	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", errEndsWithColon
		}
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		// target doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

var customAuthorityDialler = func(authority string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, authority)
	}
}

var customAuthorityResolver = func(authority string) (netResolver, error) {
	host, port, err := parseTarget(authority, defaultDNSSvrPort)
	if err != nil {
		return nil, err
	}

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial:     customAuthorityDialler(authorityWithPort),
	}, nil
}
