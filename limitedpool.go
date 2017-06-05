package gopool

import (
	"runtime"
	"sync"
)

type limitedPool struct {
	new     func() interface{}
	release func(interface{})
	check   func(interface{}) bool
	queue   chan interface{}
	max     uint
	current uint
	mu      sync.Mutex
}

func NewLPool(initial, max uint, new func() interface{}, release func(interface{}), check func(interface{}) bool) (Pool, error) {
	if initial > max || new == nil {
		return nil, ErrorInvalidParameters
	}

	pool := &limitedPool{
		queue:   make(chan interface{}, max),
		new:     new,
		release: release,
		check:   check,
		max:     max,
		current: initial,
	}

	runtime.SetFinalizer(pool, poolFilanizer)

	for ; initial > 0; initial-- {
		pool.queue <- pool.new()
	}

	return pool, nil
}

func (pool *limitedPool) Get() interface{} {
	if pool.queue == nil {
		// pool is aleardy destroyed, return nothing
		return nil
	}

	for {
		select {
		case item := <-pool.queue:
			if pool.check != nil && pool.check(item) == false {
				if pool.release != nil {
					pool.release(item)
				}
				item = pool.new()
			}
			return item
		default:
			pool.mu.Lock()
			if pool.current < pool.max {
				defer pool.mu.Unlock()
				pool.current++
				return pool.new()
			}
			pool.mu.Unlock()
			// wait for released item
			if item, ok := <-pool.queue; ok {
				return item
			}
			// nothing to return
			return nil
		}
	}
}

func (pool *limitedPool) Put(item interface{}) {
	select {
	case pool.queue <- item:
		return
	default:
		// if pool is full or destroyed, just release item
		pool.mu.Lock()
		defer pool.mu.Unlock()
		if pool.current > 0 {
			pool.current--
		} else {
			// panic! we are releasing more than allocated objects
		}
		if pool.release != nil {
			pool.release(item)
		}
		return
	}
}

func (pool *limitedPool) destroy() {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if pool.queue == nil {
		// pool is aleardy destroyed
		return
	}
	close(pool.queue)
	for item := range pool.queue {
		if pool.release != nil {
			pool.release(item)
		}
	}
	pool.queue = nil
	pool.current = 0
	pool.max = 0
}
