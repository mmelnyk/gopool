# gopool
gopool is very basic implementation for controlled pool

For our projects we needed pool with controlled pool (e.g. pool of connections) with possibility to manage pool size and the way of reallocating and destroying items.

TODO: add more info

Pool interface is compatible with sync/pool implementation.

Example:
```go
	...
	pool, err := NewLPool(2, 5, func() interface{} {
		fmt.Println("Shared object allocated")
		return &t{}
	}, func(i interface{}) {
		fmt.Println("Shared released")
	}, func(i interface{}) bool {
		fmt.Println("Validate shared object")
		return true
	})

	if err!=nil {
		panic("Error during pool creating:"+err)
	}

	for j := 1; j < 10; j++ {
		go func() {
			if v1,ok := pool.Get().(*t); ok {
				defer pool.Put(v1)
				v1.print()
			}
		}()
	}
	...
```