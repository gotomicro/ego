package ali

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const requestIDHeader = "x-log-requestid"

const (
	httpScheme  = "http://"
	httpsScheme = "https://"
)

type badResError struct {
	body   string
	header map[string][]string
	code   int
}

func (e badResError) String() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

func (e badResError) Error() string {
	return e.String()
}

func newBadResError(body []byte, header map[string][]string, httpCode int) *badResError {
	return &badResError{
		body:   string(body),
		header: header,
		code:   httpCode,
	}
}

type cliError struct {
	HTTPCode  int    `json:"httpCode"`
	Code      string `json:"errorCode"`
	Message   string `json:"errorMessage"`
	RequestID string `json:"requestID"`
}

func (e cliError) String() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

func (e cliError) Error() string {
	return e.String()
}

// newClientError new client error
func newClientError(err error) *cliError {
	if err == nil {
		return nil
	}
	if clientError, ok := err.(*cliError); ok {
		return clientError
	}
	clientError := new(cliError)
	clientError.HTTPCode = -1
	clientError.Code = "ClientError"
	clientError.Message = err.Error()
	return clientError
}

// request sends a request to alibaba cloud Log Service.
// @note if error is nil, you must call http.Response.Body.Close() to finalize reader
func (p *logProject) request(method, uri string, headers map[string]string, body []byte) (*resty.Response, error) {
	// The caller should provide 'x-log-bodyrawsize' header
	if _, ok := headers["x-log-bodyrawsize"]; !ok {
		return nil, newClientError(fmt.Errorf("can't find 'x-log-bodyrawsize' header"))
	}
	// headers["Host"] = strings.Trim(p.Host, "http://")
	headers["Date"] = nowRFC1123()
	headers["x-log-apiversion"] = version
	headers["x-log-signaturemethod"] = signatureMethod
	headers["User-Agent"] = "ego-sdk"

	// Access with token
	if p.securityToken != "" {
		headers["x-acs-security-token"] = p.securityToken
	}

	if body != nil {
		bodyMD5 := fmt.Sprintf("%X", md5.Sum(body))
		headers["Content-MD5"] = bodyMD5
		if _, ok := headers["Content-Type"]; !ok {
			return nil, newClientError(fmt.Errorf("can't find 'Content-Type' header"))
		}
	}

	// Calc Authorization
	// Authorization = "SLS <AccessKeyId>:<Signature>"
	digest, err := signature(p.accessKeySecret, method, uri, headers)
	if err != nil {
		return nil, newClientError(err)
	}
	auth := fmt.Sprintf("SLS %v:%v", p.accessKeyID, digest)
	headers["Authorization"] = auth

	res, err := p.cli.R().SetHeaders(headers).SetBody(body).Execute(method, uri)
	if err != nil {
		return nil, err
	}

	// Parse the ali error from body.
	if res.StatusCode() != http.StatusOK {
		err := &cliError{}
		err.HTTPCode = res.StatusCode()
		if jErr := json.Unmarshal(res.Body(), err); jErr != nil {
			return nil, newBadResError(res.Body(), res.Header(), res.StatusCode())
		}
		err.RequestID = res.Header().Get(requestIDHeader)
		return nil, err
	}
	return res, nil
}
