package xcycle

import (
	"fmt"
	"testing"
	"time"
)

func TestCycleDone(t *testing.T) {
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond)
		return nil
	})
	go func() {
		select {
		case <-c.Done():
			state = "done"
		case <-time.After(time.Second):
			state = "close"
		}
		c.Close()
	}()
	<-c.Wait(false)
	want := "done"
	if state != want {
		t.Errorf("TestCycleDone error want: %v, ret: %v\r\n", want, state)
	}
}

func TestCycleClose(t *testing.T) {
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Millisecond * 100)
		return nil
	})
	go func() {
		select {
		case <-c.Done():
			state = "done"
		case <-time.After(time.Millisecond):
			state = "close"
		}
		c.Close()
	}()
	<-c.Wait(false)
	want := "close"
	if state != want {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want, state)
	}
}

func TestCycleDoneAndClose(t *testing.T) {
	ch := make(chan string, 2)
	state := ""
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond * 100)
		return nil
	})
	go func() {
		c.DoneAndClose()
		ch <- "close"
	}()
	<-c.Wait(false)
	want := "close"
	state = <-ch
	if state != want {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want, state)
	}
}
func TestCycleWithError(t *testing.T) {
	c := NewCycle()
	c.Run(func() error {
		return fmt.Errorf("run error")
	})
	err := <-c.Wait(false)
	want := fmt.Errorf("run error")
	if err.Error() != want.Error() {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want.Error(), err.Error())
	}
}
