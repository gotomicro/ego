package ejob

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"

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

func (s *store) cloneCache() map[string]*Component {
	s.RLock()
	defer s.RUnlock()
	res := make(map[string]*Component)
	for jobName, component := range s.cache {
		res[jobName] = component
	}
	return res
}

// Container defines a component instance.
type Container struct {
	config *Config
	logger *elog.Component
}

// DefaultContainer returns an default container.
// Deprecated Use ejob.Job()
func DefaultContainer() *Container {
	return &Container{
		config: defaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Build constructs a specific component from container.
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
		w.Header().Set("X-Ego-Job-Err", "jobName not exist")
		w.WriteHeader(400)
		return
	}
	w.Header().Set("X-Ego-Job-Name", jobName)

	//
	jobRunID := r.Header.Get("X-Ego-Job-RunID")
	if jobRunID == "" {
		w.Header().Set("X-Ego-Job-Err", fmt.Sprintf("jobName: %s, jobRunID not exist", jobName))
		w.WriteHeader(400)
		return
	}
	w.Header().Set("X-Ego-Job-RunID", jobRunID)

	var comp *Component
	storeCache.RLock()
	comp, ok := storeCache.cache[jobName]
	storeCache.RUnlock()
	if !ok {
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
	// w.WriteHeader(200)
}

// HandleJobList job列表
func HandleJobList(w http.ResponseWriter, r *http.Request) {
	jobMap := storeCache.cloneCache()
	jobList := make([]string, 0, len(jobMap))
	for jobName := range jobMap {
		jobList = append(jobList, jobName)
	}
	buf, err := json.Marshal(jobList)
	if err != nil {
		elog.Error("HandleJobList json.Marshal failed", zap.Error(err), zap.Any("jobList", jobList))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(buf)
	if err != nil {
		elog.Error("HandleJobList write failed", zap.Error(err))
	}
}
