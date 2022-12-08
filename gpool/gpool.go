package gpool

import (
	"fmt"
	"sync"
)

type GPool struct {
	lock      sync.Mutex
	maxCount  int64
	curCount  int64
	waitGroup sync.WaitGroup
	jobs      chan interface{}
	fun       func(param interface{})
	closeCh   chan bool
}

var NilErr = fmt.Errorf("param can not be nil")
var CloseErr = fmt.Errorf("pool was closed")

func NewGPool(maxCount int64, fun func(param interface{})) *GPool {
	return &GPool{
		maxCount: maxCount,
		fun:      fun,
		jobs:     make(chan interface{}, 10),
		closeCh:  make(chan bool),
	}
}

func (g *GPool) Run(param interface{}) error {
	if param == nil {
		return NilErr
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	select {
	case <-g.closeCh:
		return CloseErr
	default:
		if g.curCount < g.maxCount {
			g.waitGroup.Add(1)
			g.curCount++
			go func() {
				defer func() {
					if p := recover(); p != nil {
						fmt.Printf("%#v\n", p)
					}
					g.waitGroup.Done()
				}()
				// consumer
				g.worker()
			}()
		}
		// producer
		g.jobs <- param
		return nil
	}
}

func (g *GPool) Clear() {
	g.lock.Lock()
	defer g.lock.Unlock()
	select {
	case <-g.closeCh:
	default:
		for g.curCount > 0 {
			g.curCount--
			g.jobs <- nil
		}
	}
}

func (g *GPool) Close() {
	g.lock.Lock()
	select {
	case <-g.closeCh:
	default:
		close(g.closeCh)
		close(g.jobs)
	}
	g.lock.Unlock()
	g.waitGroup.Wait()
}

func (g *GPool) worker() {
	for j := range g.jobs {
		if j == nil {
			break
		}
		g.fun(j)
	}
}
