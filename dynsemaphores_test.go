package dynsemaphore

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestDynSemaphoreGoRoutines(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		expected    int
	}{
		{
			name:        "test 2",
			concurrency: 2,
			expected:    2,
		},
		{
			name:        "test 0",
			concurrency: 0,
			expected:    99999,
		},
		{
			name:        "test cpu",
			concurrency: -1,
			expected:    runtime.NumCPU(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			d := New(tt.concurrency)
			for i := 0; i < runtime.NumCPU()*10; i++ {
				wg.Add(1)
				go func(d *DynSemaphore, i int) {
					d.Access()
					time.Sleep(5 * time.Millisecond)
					if d.MaxConcurrency > 0 && d.n > d.MaxConcurrency {
						t.Errorf("Concurrency exceeded: %d > %d\n", d.n, tt.expected)
					}
					d.Release()
					wg.Done()
				}(d, i)
				time.Sleep(time.Millisecond)
			}
			wg.Wait()
		})
	}
}
