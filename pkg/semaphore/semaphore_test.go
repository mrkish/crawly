package semaphore_test

import (
	"crawly/pkg/semaphore"
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_Weighted(t *testing.T) {
	tt := []struct {
		name       string
		workers    int
		iterations int
	}{
		{
			name:       "basic semaphore works as expected",
			workers:    3,
			iterations: 5,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sem := semaphore.New(tc.workers)
			wg := sync.WaitGroup{}
			for i := range tc.iterations {
				fmt.Printf("acquiring sem: test loop: %d\n", i)
				sem.Acquire()
				wg.Add(1)
				go func() {
					defer sem.Free()
					defer wg.Done()
					fmt.Printf("test loop: %d\n", i)
					time.Sleep(time.Second * 1)
				}()
			}
			wg.Wait()
		})
	}
}
