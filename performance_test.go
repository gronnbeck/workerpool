package workerpool

import (
	"context"
	"testing"
	"time"
)

func Test_smallWorksmallPool(t *testing.T) {
	runPerformanceTest(t, 1000, 1000, nil, 10*time.Second)
}

func Test_mediumWorksmallPool(t *testing.T) {
	runPerformanceTest(t, 1000, 100000, nil, 10*time.Second)
}

func Test_largeWorkLargePool(t *testing.T) {
	runPerformanceTest(t, 10000, 1000000, nil, 60*time.Second)
}

func runPerformanceTest(t *testing.T, poolSize, workSize int, work func() error, timeout time.Duration) {
	pool := New(poolSize)
	pool.Start()
	defer pool.Stop()
	res := make(chan error)

	if work == nil {
		work = func() error {
			return nil
		}
	}

	go func() {
		for i := 0; i < workSize; i++ {
			ctx := context.Background()
			pool.Queue(ctx, func() {
				res <- work() // some function that does work
			})
		}
	}()

	crashAfter := time.After(timeout)
	for i := 0; i < workSize; i++ {
		select {
		case <-crashAfter:
			t.Log("Test timeout")
			t.Fail()
		case <-time.After(50 * time.Millisecond):
			t.Log("Test timed out")
			t.Fail()
		case err := <-res:
			if err != nil {
			}
		}
	}

}
