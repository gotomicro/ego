package ali

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/proto"
	"github.com/pierrec/lz4"

	"github.com/gotomicro/ego/core/elog/ali/pb"
)

const (
	version         = "0.6.0"     // SDK version
	signatureMethod = "hmac-sha1" // Signature method
)

// errorMessage message in SLS HTTP response.
type errorMessage struct {
	Code    string `json:"errorCode"`
	Message string `json:"errorMessage"`
}

// LogProject defines the Ali project detail
type LogProject struct {
	name            string // project name
	endpoint        string // IP or hostname of SLS endpoint
	accessKeyID     string
	accessKeySecret string
	host            string
	usingHTTP       bool
	securityToken   string
	cli             *resty.Client
}

// LogStore stores the logs
type LogStore struct {
	Name           string `json:"logstoreName"`
	TTL            int
	ShardCount     int
	CreateTime     uint32
	LastModifyTime uint32
	project        *LogProject
}

// ListLogStore returns all logstore names of project p.
func (p *LogProject) ListLogStore() (storeNames []string, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}

	uri := fmt.Sprintf("/logstores")
	r, err := p.request("GET", uri, h, nil)
	if err != nil {
		return
	}

	if r.StatusCode() != http.StatusOK {
		errMsg := &errorMessage{}
		if err = json.Unmarshal(r.Body(), errMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal logstore response, %w", err)
		}
		return nil, fmt.Errorf("%v:%v", errMsg.Code, errMsg.Message)
	}

	type Body struct {
		Count     int
		LogStores []string
	}
	body := &Body{}
	if err = json.Unmarshal(r.Body(), body); err != nil {
		return
	}
	storeNames = body.LogStores
	return
}

func (p *LogProject) parseEndpoint() {
	scheme := httpScheme // default to http scheme
	host := p.endpoint

	if strings.HasPrefix(p.endpoint, httpScheme) {
		host = strings.TrimPrefix(p.endpoint, scheme)
	} else if strings.HasPrefix(p.endpoint, httpsScheme) {
		scheme = httpsScheme
		host = strings.TrimPrefix(p.endpoint, scheme)
	}

	if ipRegex.MatchString(host) {
		// use direct ip proxy
		u, err := url.Parse(fmt.Sprintf("%s%s", scheme, host))
		if err != nil {
			return
		}
		cli := p.cli.GetClient()
		cli.Transport = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		p.cli = resty.NewWithClient(cli)
	}
	if p.name == "" {
		p.host = fmt.Sprintf("%s%s", scheme, host)
	} else {
		p.host = fmt.Sprintf("%s%s.%s", scheme, p.name, host)
	}
}

// GetLogStore returns logstore according by logstore name.
func (p *LogProject) GetLogStore(name string) (s *LogStore, err error) {
	h := map[string]string{"x-log-bodyrawsize": "0"}
	r, err := p.request("GET", "/logstores/"+name, h, nil)
	if err != nil {
		return
	}

	if r.StatusCode() != http.StatusOK {
		errMsg := &errorMessage{}
		if err = json.Unmarshal(r.Body(), errMsg); err != nil {
			return nil, fmt.Errorf("failed to get logstore")
		}
		return nil, fmt.Errorf("%v:%v", errMsg.Code, errMsg.Message)
	}

	s = &LogStore{}
	if err = json.Unmarshal(r.Body(), s); err != nil {
		return
	}
	s.project = p
	return
}

// PutLogs puts logs into logstore.
// The callers should transform user logs into LogGroup.
func (s *LogStore) PutLogs(lg *pb.LogGroup) (err error) {
	body, err := proto.Marshal(lg)
	if err != nil {
		return
	}

	// Compresse body with lz4
	out := make([]byte, lz4.CompressBlockBound(len(body)))
	n, err := lz4.CompressBlock(body, out, nil)
	if err != nil {
		return
	}

	h := map[string]string{
		"x-log-compresstype": "lz4",
		"x-log-bodyrawsize":  strconv.Itoa(len(body)),
		"Content-Type":       "application/x-protobuf",
	}

	uri := fmt.Sprintf("/logstores/%v", s.Name)
	r, err := s.project.request("POST", uri, h, out[:n])
	if err != nil {
		return
	}

	if r.StatusCode() != http.StatusOK {
		errMsg := &errorMessage{}
		err = json.Unmarshal(r.Body(), errMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal putLogs response, %w", err)
		}
		return fmt.Errorf("%v:%v", errMsg.Code, errMsg.Message)
	}
	return
}
