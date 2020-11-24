package eetcd

import (
	"context"
	"github.com/gotomicro/ego/core/elog"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/util/xgo"
)

// Watch A watch only tells the latest revision
type Watch struct {
	revision     int64
	cancel       context.CancelFunc
	eventChan    chan *clientv3.Event
	lock         *sync.RWMutex
	incipientKVs []*mvccpb.KeyValue
}

// C ...
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

// IncipientKeyValues incipient key and values
func (w *Watch) IncipientKeyValues() []*mvccpb.KeyValue {
	return w.incipientKVs
}

// NewWatch ...
func (client *Component) WatchPrefix(ctx context.Context, prefix string) (*Watch, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var w = &Watch{
		revision:     resp.Header.Revision,
		eventChan:    make(chan *clientv3.Event, 100),
		incipientKVs: resp.Kvs,
	}

	xgo.Go(func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.Client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					client.logger.Error(ecode.MsgWatchRequestErr, elog.FieldErrKind(ecode.ErrKindRegisterErr), elog.FieldErr(err), elog.FieldAddr(prefix))
					continue
				}
				for _, ev := range n.Events {
					select {
					case w.eventChan <- ev:
					default:
						client.logger.Error("watch etcd with prefix", elog.Any("err", "block event chan, drop event message"))
					}
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel = cancel
			if w.revision > 0 {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
			} else {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	})

	return w, nil
}

// Close close watch
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
