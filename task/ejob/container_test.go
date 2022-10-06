package ejob

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	fc := func(ctx Context) error {
		return nil
	}
	Job("test", fc)
	storeCache.RLock()
	assert.Equal(t, 1, len(storeCache.cache))
	storeCache.RUnlock()
}
func TestHandle(t *testing.T) {

	fc := func(ctx Context) error {
		return nil
	}
	Job("test", fc)
	storeCache.RLock()
	assert.Equal(t, 1, len(storeCache.cache))
	storeCache.RUnlock()
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", Handle)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	{
		resp := performRequest(mux, "POST", "/jobs")
		assert.Equal(t, 400, resp.Code)
		assert.Equal(t, "jobName not exist", resp.Header().Get("X-Ego-Job-Err"))
	}
	{
		resp := performRequest(mux, "POST", "/jobs", header{
			Key:   "X-Ego-Job-Name",
			Value: "test",
		})
		assert.Equal(t, 400, resp.Code)
		assert.Equal(t, "jobName: test, jobRunID not exist", resp.Header().Get("X-Ego-Job-Err"))
	}
	{
		resp := performRequest(mux, "POST", "/jobs", header{
			Key:   "X-Ego-Job-Name",
			Value: "test1",
		}, header{
			Key:   "X-Ego-Job-RunID",
			Value: "123",
		})
		assert.Equal(t, 400, resp.Code)
		assert.Equal(t, "job:test1 not exist", resp.Header().Get("X-Ego-Job-Err"))
	}

	{
		resp := performRequest(mux, "POST", "/jobs", header{
			Key:   "X-Ego-Job-Name",
			Value: "test",
		}, header{
			Key:   "X-Ego-Job-RunID",
			Value: "123",
		})
		assert.Equal(t, 200, resp.Code)
	}

}

func TestHandleJobList(t *testing.T) {
	fc := func(ctx Context) error {
		return nil
	}
	Job("test", fc)
	mux := http.NewServeMux()
	mux.HandleFunc("/jobList", HandleJobList)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp := performRequest(mux, "GET", "/jobList")
	jobList := make([]string, 0)
	bytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, &jobList)
	assert.Nil(t, err)
	assert.Equal(t, "test", jobList[0])
}

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
