package ali

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog/ali/pb"
	"github.com/gotomicro/ego/core/util/xcast"
)

const (
	// flushSize sets the flush size
	flushSize int = 128
)

type LogContent = pb.Log_Content

// config is the config for Ali Log
type config struct {
	encoder             zapcore.Encoder
	project             string
	endpoint            string
	accessKeyID         string
	accessKeySecret     string
	logstore            string
	topics              []string
	source              string
	flushSize           int
	flushBufferSize     int32
	flushBufferInterval time.Duration
	levelEnabler        zapcore.LevelEnabler
	apiBulkSize         int
	apiTimeout          time.Duration
	apiRetryCount       int
	apiRetryWaitTime    time.Duration
	apiRetryMaxWaitTime time.Duration
}

// writer implements LoggerInterface.
// Writes messages in keep-live tcp connection.
type writer struct {
	store      *LogStore
	group      []*pb.LogGroup
	withMap    bool
	ch         chan *pb.Log
	curBufSize *int32
	cancel     context.CancelFunc
	config
}

func retryCondition(r *resty.Response, err error) bool {
	code := r.StatusCode()
	if code == 500 || code == 502 || code == 503 {
		return true
	}
	return false
}

// newWriter creates a new ali writer
func newWriter(c config) (*writer, error) {
	w := &writer{config: c, ch: make(chan *pb.Log, c.apiBulkSize), curBufSize: new(int32)}
	p := &LogProject{
		name:            w.project,
		endpoint:        w.endpoint,
		accessKeyID:     w.accessKeyID,
		accessKeySecret: w.accessKeySecret,
	}
	p.parseEndpoint()
	p.cli = resty.New().
		SetDebug(eapp.IsDevelopmentMode()).
		SetHostURL(p.host).
		SetTimeout(c.apiTimeout).
		SetRetryCount(c.apiRetryCount).
		SetRetryWaitTime(c.apiRetryWaitTime).
		SetRetryMaxWaitTime(c.apiRetryMaxWaitTime).
		AddRetryCondition(retryCondition)
	store, err := p.GetLogStore(w.logstore)
	if err != nil {
		return nil, fmt.Errorf("getlogstroe fail,%w", err)
	}
	w.store = store

	// Create default Log Group
	w.group = append(w.group, &pb.LogGroup{
		Topic:  proto.String(""),
		Source: proto.String(w.source),
		Logs:   make([]*pb.Log, 0, w.apiBulkSize),
	})

	// Create other Log Group
	for _, topic := range w.topics {
		lg := &pb.LogGroup{
			Topic:  proto.String(topic),
			Source: proto.String(w.source),
			Logs:   make([]*pb.Log, 0, w.apiBulkSize),
		}
		w.group = append(w.group, lg)
	}

	w.Sync()
	return w, nil
}

func genLog(fields map[string]interface{}) *pb.Log {
	l := &pb.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: make([]*LogContent, 0, len(fields)),
	}
	for k, v := range fields {
		l.Contents = append(l.Contents, &LogContent{
			Key:   proto.String(k),
			Value: proto.String(xcast.ToString(v)),
		})
	}
	return l
}

func (w *writer) write(fields map[string]interface{}) (err error) {
	l := genLog(fields)
	// if bufferSize bigger then defaultBufferSize or channel is full, then flush logs
	w.ch <- l
	atomic.AddInt32(w.curBufSize, int32(l.XXX_Size()))
	if atomic.LoadInt32(w.curBufSize) >= w.flushBufferSize || len(w.ch) >= cap(w.ch) {
		err = w.flush()
	}
	return
}

func (w *writer) flush() error {
	// TODO sync pool
	var lg = *w.group[0]
	chlen := len(w.ch)
	if chlen == 0 {
		return nil
	}
	var logs = make([]*pb.Log, 0, chlen)
	logs = append(logs, <-w.ch)
L1:
	for {
		select {
		case l := <-w.ch:
			logs = append(logs, l)
		default:
			break L1
		}
	}
	lg.Logs = logs
	return w.store.PutLogs(&lg)
}

func (w *writer) Sync() {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	ticker := time.NewTicker(w.flushBufferInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := w.flush(); err != nil {
					log.Printf("writer flush fail, %s\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return
}
