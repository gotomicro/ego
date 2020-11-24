package xgo

import (
	"errors"
	"testing"
	"time"
)

var (
	fn1     = func() error { return nil }
	fn2     = func() error { return errors.New("BOOM") }
	timeout = time.After(2 * time.Second)
)

func TestParallel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "test run",
			function: testRun,
		},
		{
			scenario: "test run limit",
			function: testRunLimit,
		},
		{
			scenario: "test run limit with negative concurrency value",
			function: testRunLimitWithNegativeConcurrencyValue,
		},
		{
			scenario: "test run limit with concurrency value greater than passed functions",
			function: testRunLimitWithConcurrencyGreaterThanPassedFunctions,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testRun(t *testing.T) {
	var count int
	err := ParallelWithErrorChan(fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimit(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(2, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimitWithNegativeConcurrencyValue(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(-1, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimitWithConcurrencyGreaterThanPassedFunctions(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(3, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}
