package ejob

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gotomicro/ego/core/elog"
)

var storeCache *store

type store struct {
	sync.RWMutex
	cache map[string]*Component
}

func init() {
	storeCache = &store{
		RWMutex: sync.RWMutex{},
		cache:   make(map[string]*Component),
	}
}

// Container 容器
type Container struct {
	config *Config
	logger *elog.Component
}

// DefaultContainer 默认容器
// Deprecated Use ejob.Job()
func DefaultContainer() *Container {
	return &Container{
		config: defaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Build 构建组件
// Deprecated Use ejob.Job()
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config.Name, c.config, c.logger)
}

// Job Job方法
func Job(name string, startFunc func(Context) error) *Component {
	container := &Container{
		config: defaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
	comp := container.Build(
		WithName(name),
		WithStartFunc(startFunc),
	)
	storeCache.Lock()
	storeCache.cache[name] = comp
	storeCache.Unlock()
	return comp
}

// Handle ...
func Handle(w http.ResponseWriter, r *http.Request) {
	jobName := r.Header.Get("X-Ego-Job-Name")
	if jobName == "" {
		return
	}
	w.Header().Set("X-Ego-Job-Name", jobName)

	//
	jobRunID := r.Header.Get("X-Ego-Job-RunID")
	if jobName == "" {
		return
	}
	w.Header().Set("X-Ego-Job-RunID", jobRunID)

	var comp *Component
	storeCache.RLock()
	comp = storeCache.cache[jobName]
	storeCache.RUnlock()
	if comp == nil {
		w.Header().Set("X-Ego-Job-Err", fmt.Sprintf("job:%s not exist", jobName))
		w.WriteHeader(400)
		return
	}
	err := comp.StartHTTP(w, r)
	if err != nil {
		w.Header().Set("X-Ego-Job-Err", err.Error())
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}
