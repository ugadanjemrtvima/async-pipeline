package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestByIlia(t *testing.T) {

	var received uint32
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- uint32(1)
			out <- uint32(3)
			out <- uint32(4)
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				out <- val.(uint32) * 3
				time.Sleep(time.Millisecond * 100)
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("collected", val)
				atomic.AddUint32(&received, val.(uint32))
			}
		}),
	}

	start := time.Now()

	ExecutePipeline(freeFlowJobs...)

	end := time.Since(start)

	expectedTime := time.Millisecond * 350

	if end > expectedTime {
		t.Errorf("execution too long\nGot: %s\nExpected: <%s", end, expectedTime)
	}

	if received != (1+3+4)*3 {
		t.Errorf("f3 have not collected inputs, received = %d", received)
	}
}
