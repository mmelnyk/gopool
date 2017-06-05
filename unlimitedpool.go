package gopool

import (
	"runtime"
	"sync"
)

type unlimitedPool struct {
	new     func() interface{}
	release func(interface{})
	check   func(interface{}) bool
	queue   chan interface{}
	mu      sync.Mutex
}

func NewPool(initial, max uint, new func() interface{}, release func(interface{}), check func(interface{}) bool) (Pool, error) {
	if initial > max || new == nil {
		return nil, ErrorInvalidParameters
	}

	pool := &unlimitedPool{
		queue:   make(chan interface{}, max),
		new:     new,
		release: release,
		check:   check,
	}

	runtime.SetFinalizer(pool, poolFilanizer)

	for ; initial > 0; initial-- {
		pool.queue <- pool.new()
	}

	return pool, nil
}

func (pool *unlimitedPool) Get() interface{} {
	if pool.queue == nil {
		// pool aleardy destroyed, return nothing
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
			return pool.new()
		}
	}
}

func (pool *unlimitedPool) Put(item interface{}) {
	select {
	case pool.queue <- item:
		return
	default:
		// pool is full or destroyed, destroy item
		if pool.release != nil {
			pool.release(item)
		}
		return
	}
}

func (pool *unlimitedPool) destroy() {
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
}
