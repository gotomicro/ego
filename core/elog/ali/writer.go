package ali

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/publicsuffix"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog/ali/pb"
	"github.com/gotomicro/ego/core/emetric"
)

const (
	// observe interval
	observeInterval = 5 * time.Second
	// apiBulkMinSize sets bulk minimal size
	apiBulkMinSize = 256
)

// logContent ...
type logContent = pb.Log_Content

// config is the config for Ali Log
type config struct {
	encoder                zapcore.Encoder
	project                string
	endpoint               string
	accessKeyID            string
	accessKeySecret        string
	logstore               string
	maxQueueSize           int
	flushBufferSize        int32
	flushBufferInterval    time.Duration
	levelEnabler           zapcore.LevelEnabler
	apiBulkSize            int
	apiTimeout             time.Duration
	apiRetryCount          int
	apiRetryWaitTime       time.Duration
	apiRetryMaxWaitTime    time.Duration
	apiMaxIdleConns        int
	apiIdleConnTimeout     time.Duration
	apiMaxIdleConnsPerHost int
	fallbackCore           zapcore.Core
}

// writer implements LoggerInterface.
// Writes messages in keep-live tcp connection.
type writer struct {
	fallbackCore zapcore.Core
	store        *logStore
	ch           chan *pb.Log
	lock         sync.Mutex
	curBufSize   *int32
	cancel       context.CancelFunc
	config
}

func retryCondition(r *resty.Response, err error) bool {
	return r.StatusCode() != 200
}

// newWriter creates a new ali writer
func newWriter(c config) (*writer, error) {
	if c.apiBulkSize >= c.maxQueueSize {
		c.apiBulkSize = c.maxQueueSize
	}
	if c.apiBulkSize < apiBulkMinSize {
		c.apiBulkSize = apiBulkMinSize
	}
	w := &writer{config: c, ch: make(chan *pb.Log, c.maxQueueSize), curBufSize: new(int32)}
	p := &logProject{
		name:            w.project,
		endpoint:        w.endpoint,
		accessKeyID:     w.accessKeyID,
		accessKeySecret: w.accessKeySecret,
	}
	p.initHost()
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	p.cli = resty.NewWithClient(&http.Client{
		Transport: createTransport(c),
		Jar:       cookieJar,
	}).
		SetDebug(eapp.IsDevelopmentMode()).
		SetHostURL(p.host).
		SetTimeout(c.apiTimeout).
		SetRetryCount(c.apiRetryCount).
		SetRetryWaitTime(c.apiRetryWaitTime).
		SetRetryMaxWaitTime(c.apiRetryMaxWaitTime).
		AddRetryCondition(retryCondition)
	store, err := p.getLogStore(w.logstore)
	if err != nil {
		return nil, fmt.Errorf("getlogstroe fail,%w", err)
	}
	w.store = store
	w.fallbackCore = c.fallbackCore
	w.sync()
	w.observe()
	return w, nil
}

func createTransport(c config) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	if c.apiMaxIdleConnsPerHost == 0 {
		c.apiMaxIdleConnsPerHost = c.apiMaxIdleConns
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          c.apiMaxIdleConns,
		IdleConnTimeout:       c.apiIdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   c.apiMaxIdleConnsPerHost,
	}
}

func genLog(fields map[string]interface{}) *pb.Log {
	l := &pb.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: make([]*logContent, 0, len(fields)),
	}
	for k, v := range fields {
		valStr, err := toStringE(v)
		if err != nil {
			log.Printf("toString fail, %s", err.Error())
		}
		l.Contents = append(l.Contents, &logContent{
			Key:   proto.String(k),
			Value: proto.String(valStr),
		})
	}
	return l
}

type jsonEncoder struct {
	enc *json.Encoder
}

var jsonEncPool = sync.Pool{New: func() interface{} {
	return &jsonEncoder{}
}}

var bufPool = buffer.NewPool()

func objToString(obj interface{}) (string, error) {
	enc := jsonEncPool.Get().(*jsonEncoder)
	buf := bufPool.Get()
	enc.enc = json.NewEncoder(buf)

	defer func() {
		buf.Reset()
		buf.Free()
		enc.enc = nil
		jsonEncPool.Put(enc)
	}()

	enc.enc.SetEscapeHTML(false)
	if err := enc.enc.Encode(obj); err != nil {
		return "", err
	}
	buf.TrimNewline()
	return buf.String(), nil
}

func (w *writer) write(fields map[string]interface{}) (err error) {
	l := genLog(fields)
	// if bufferSize bigger then defaultBufferSize or channel is full, then flush logs
	w.ch <- l
	atomic.AddInt32(w.curBufSize, int32(l.XXX_Size()))
	if atomic.LoadInt32(w.curBufSize) >= w.flushBufferSize || len(w.ch) >= cap(w.ch) {
		err = w.flush()
		atomic.StoreInt32(w.curBufSize, 0)
	}
	return
}

func (w *writer) flush() error {
	w.lock.Lock()
	entriesChLen := len(w.ch)
	if entriesChLen == 0 {
		w.lock.Unlock()
		return nil
	}
	var waitedEntries = make([]*pb.Log, 0, entriesChLen)
	for i := 0; i < entriesChLen; i++ {
		waitedEntries = append(waitedEntries, <-w.ch)
	}
	w.lock.Unlock()

	chunks := int(math.Ceil(float64(len(waitedEntries)) / float64(w.apiBulkSize)))
	wg := sync.WaitGroup{}
	wg.Add(chunks)
	for i := 0; i < chunks; i++ {
		go func(start int) {
			end := (start + 1) * w.apiBulkSize
			if end > len(waitedEntries) {
				end = len(waitedEntries)
			}
			lg := pb.LogGroup{Logs: waitedEntries[start:end]}
			if e := w.store.putLogs(&lg); e != nil {
				log.Println("[sls] putLogs to sls fail,try to write to fallback logger now,", e)
				// if error occurs we put logs to fallback logger
				w.writeToFallbackLogger(lg)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	return nil
}

func (w *writer) writeToFallbackLogger(lg pb.LogGroup) {
	for _, v := range lg.Logs {
		fields := make([]zapcore.Field, len(v.Contents))
		for i, val := range v.Contents {
			fields[i] = zap.String(val.GetKey(), val.GetValue())
		}
		if e := w.fallbackCore.Write(zapcore.Entry{Time: time.Now()}, fields); e != nil {
			log.Println("[sls] fallbackCore write fail,", e)
		}
	}
}

func (w *writer) sync() {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	ticker := time.NewTicker(w.flushBufferInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := w.flush(); err != nil {
					log.Printf("[sls] writer flush fail, %s\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *writer) observe() {
	go func() {
		for {
			emetric.LibHandleSummary.Observe(float64(len(w.ch)), "elog", "ali_waited_entries")
			time.Sleep(observeInterval)
		}
	}()
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// ToStringE casts an interface to a string type.
func toStringE(i interface{}) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "null", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return objToString(i)
	}
}
