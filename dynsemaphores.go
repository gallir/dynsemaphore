package dynsemaphore

import (
	"runtime"
	"sync"
)

type DynSemaphore struct {
	MaxConcurrency int
	n              int
	mu             sync.Mutex
	cond           *sync.Cond
}

func New(c int) *DynSemaphore {
	if c < 0 {
		c = runtime.NumCPU()
	}
	d := &DynSemaphore{
		MaxConcurrency: c,
	}
	d.cond = sync.NewCond(&d.mu)
	return d
}

func (d *DynSemaphore) Access() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for d.MaxConcurrency > 0 && d.n >= d.MaxConcurrency {
		d.cond.Wait()
	}
	d.n++
}

func (d *DynSemaphore) Release() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.n--
	d.cond.Signal()
}

func (d *DynSemaphore) SetConcurrency(c int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if c < 0 {
		c = runtime.NumCPU()
	}
	d.MaxConcurrency = c
}
