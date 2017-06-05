package gopool

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Pool provides generic interface for custom pool 
type Pool interface {
	Get() interface{}
	Put(interface{})
}

var (
	ErrorInvalidParameters = errors.New("Invalid Parameters")
)

func poolFilanizer(pool interface{}) {
	switch v := pool.(type) {
	case *unlimitedPool:
		v.destroy()
	case *limitedPool:
		v.destroy()
	}
}
