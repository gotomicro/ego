package egrpc

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/eregistry"
)

var hostLookupTbl = struct {
	sync.Mutex
	tbl map[string][]string
}{
	tbl: map[string][]string{
		"foo.bar.com":          {"1.2.3.4", "5.6.7.8"},
		"ipv4.single.fake":     {"1.2.3.4"},
		"srv.ipv4.single.fake": {"2.4.6.8"},
		"srv.ipv4.multi.fake":  {},
		"srv.ipv6.single.fake": {},
		"srv.ipv6.multi.fake":  {},
		"ipv4.multi.fake":      {"1.2.3.4", "5.6.7.8", "9.10.11.12"},
		"ipv6.single.fake":     {"2607:f8b0:400a:801::1001"},
		"ipv6.multi.fake":      {"2607:f8b0:400a:801::1001", "2607:f8b0:400a:801::1002", "2607:f8b0:400a:801::1003"},
	},
}

func hostLookup(host string) ([]string, error) {
	hostLookupTbl.Lock()
	defer hostLookupTbl.Unlock()
	if addrs, ok := hostLookupTbl.tbl[host]; ok {
		return addrs, nil
	}
	return nil, &net.DNSError{
		Err:         "hostLookup error",
		Name:        host,
		Server:      "fake",
		IsTemporary: true,
	}
}

const colonDefaultPort = ":" + defaultPort

func TestDNSRegistry(t *testing.T) {
	tests := []struct {
		target   string
		addrWant []resolver.Address
		addrNext []resolver.Address
	}{
		{
			"foo.bar.com",
			[]resolver.Address{{Addr: "1.2.3.4" + colonDefaultPort}, {Addr: "5.6.7.8" + colonDefaultPort}},
			[]resolver.Address{{Addr: "1.2.3.4" + colonDefaultPort}},
		},
		{
			"ipv4.single.fake",
			[]resolver.Address{{Addr: "1.2.3.4" + colonDefaultPort}},
			[]resolver.Address{{Addr: "1.2.3.4" + colonDefaultPort}},
		},
		{
			"1.3.5.7:443",
			[]resolver.Address{{Addr: "1.3.5.7" + colonDefaultPort}},
			[]resolver.Address{{Addr: "1.3.5.7" + colonDefaultPort}},
		},
	}
	nc := replaceNetFunc(make(chan struct{}))
	defer nc()

	for _, a := range tests {
		b := DNSRegistry()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		r, err := b.WatchServices(ctx, eregistry.Target{Endpoint: a.target})
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		tr, ok := b.(*dnsRegistry).r.resolver.(*testResolver)
		if !ok {
			t.Fatalf("delegate resolver returned unexpected type: %T\n", tr)
		}
		var endpoints eregistry.Endpoints
		tm := time.NewTicker(2 * time.Second)
	L:
		for {
			time.Sleep(200 * time.Millisecond)
			select {
			case endpoints = <-r:
				break L
			case <-tm.C:
				t.Fatalf("timeout")
			}
		}

		get, want := []string{}, []string{}
		for k := range endpoints.Nodes {
			get = append(get, k)
		}
		for _, v := range a.addrWant {
			want = append(want, v.Addr)
		}

		if !assert.ElementsMatch(t, get, want) {
			t.Errorf("Resolved addresses of target: %+v, want %+v\n", get, want)
		}
	}
}

type testResolver struct{}

func (tr *testResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	return hostLookup(host)
}

func (*testResolver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	return srvLookup(service, proto, name)
}

func (*testResolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	return nil, nil
}

var srvLookupTbl = struct {
	sync.Mutex
	tbl map[string][]*net.SRV
}{
	tbl: map[string][]*net.SRV{
		"_grpclb._tcp.srv.ipv4.single.fake": {&net.SRV{Target: "ipv4.single.fake", Port: 1234}},
		"_grpclb._tcp.srv.ipv4.multi.fake":  {&net.SRV{Target: "ipv4.multi.fake", Port: 1234}},
		"_grpclb._tcp.srv.ipv6.single.fake": {&net.SRV{Target: "ipv6.single.fake", Port: 1234}},
		"_grpclb._tcp.srv.ipv6.multi.fake":  {&net.SRV{Target: "ipv6.multi.fake", Port: 1234}},
	},
}

func srvLookup(service, proto, name string) (string, []*net.SRV, error) {
	cname := "_" + service + "._" + proto + "." + name
	srvLookupTbl.Lock()
	defer srvLookupTbl.Unlock()
	if srvs, cnt := srvLookupTbl.tbl[cname]; cnt {
		return cname, srvs, nil
	}
	return "", nil, &net.DNSError{
		Err:         "srvLookup error",
		Name:        cname,
		Server:      "fake",
		IsTemporary: true,
	}
}

func replaceNetFunc(ch chan struct{}) func() {
	oldResolver := defaultResolver
	defaultResolver = &testResolver{}

	return func() {
		defaultResolver = oldResolver
	}
}
