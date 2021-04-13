package ecron

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"

	"github.com/gotomicro/ego/core/econf"
)

func testBuildComp(name, config string) (c *Component, err error) {
	err = econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	if err != nil {
		return nil, err
	}
	comp := Load(name).Build()
	return comp, nil
}

func TestComponent_Name(t *testing.T) {
	name := "cron.syncXxx"
	comp, err := testBuildComp(name, `[cron.syncXxx]
spec = "0 0 1 1 *"`)
	if err != nil {
		t.Errorf("load config failed. err=%s", err.Error())
		return
	}

	if comp.Name() != name {
		t.Errorf("expect name = %s. got %s", name, comp.Name())
	}
}

func TestComponent_PackageName(t *testing.T) {
	comp, err := testBuildComp("test", `[test]
spec = "0 0 1 1 *"`)
	if err != nil {
		t.Errorf("load config failed. err=%s", err.Error())
		return
	}
	if comp.PackageName() != PackageName {
		t.Errorf("expect PackName() = %s", PackageName)
	}
}

func TestComponent_Init(t *testing.T) {
	comp, err := testBuildComp("test", `[test]
spec = "0 0 1 1 *"`)
	if err != nil {
		t.Errorf("load config failed. err=%s", err.Error())
		return
	}
	if err := comp.Init(); err != nil {
		t.Errorf("expect Init() returns nil, got %v", err)
	}
}

func TestRunJob(t *testing.T) {
	name := "test"
	config := `[test]
enableImmediatelyRun = true
spec = "0 0 1 1 *"`
	err := econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	if err != nil {
		t.Errorf("load config failed. err=%s", err.Error())
		return
	}
	invoked := 0
	comp := Load(name).Build(WithJob(func(ctx context.Context) error {
		invoked++
		return nil
	}))

	go func() {
		// wait for start
		time.Sleep(time.Second)

		err := comp.Stop()
		if err != nil {
			t.Errorf("Stop() returns err: %s", err.Error())
			return
		}
	}()

	err = comp.Start()
	if err != nil {
		t.Errorf("Start() returns err: %s", err.Error())
		return
	}

	if invoked != 1 {
		t.Errorf("expect 'invoked' = 1, got %d", invoked)
	}
}

func TestRunDistributedJob(t *testing.T) {
	mtx := sync.Mutex{}
	invoked := 0
	lastNode := ""

	job := func(key string) FuncJob {
		return func(ctx context.Context) error {
			mtx.Lock()
			defer mtx.Unlock()

			invoked++
			t.Logf("job invoked %dth", invoked)
			t.Logf("%s is running", key)
			if lastNode != "" && lastNode != key {
				t.Errorf("job running on multi nodes. lastNode=%s, thisNode=%s", lastNode, key)
			}
			lastNode = key
			return nil
		}
	}

	lock := &mockLock{}

	config := `[test]
enableSeconds = true
enableDistributedTask = true
spec = "0/1 * * * * *"
`
	err := econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	if err != nil {
		t.Errorf("load config failed. err=%s", err.Error())
		return
	}

	runCronJob := func(key string) func() error {
		return func() error {
			comp := Load("test").Build(
				WithJob(job(key)),
				WithParser(cron.NewParser(cron.Second)),
				WithLock(lock),
			)

			go func() {
				// wait for running
				time.Sleep(5 * time.Second)
				_ = comp.Stop()
			}()

			err := comp.Start()
			if err != nil {
				t.Errorf("Start() returns err: %s", err.Error())
				return err
			}

			return nil
		}
	}

	// mock 10 nodes
	cronCount := 10
	eg := errgroup.Group{}
	for i := 0; i < cronCount; i++ {
		eg.Go(runCronJob(fmt.Sprintf("node-%d", i)))
	}

	err = eg.Wait()
	if err != nil {
		t.Errorf("run cron job failed, go error: %s", err.Error())
		return
	}
}
