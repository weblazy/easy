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
}

func NewGPool(maxCount int64, fun func(param interface{})) *GPool {
	return &GPool{
		maxCount: maxCount,
		fun:      fun,
		jobs:     make(chan interface{}, 10),
	}
}

func (g *GPool) Run(param interface{}) {
	if g.curCount < g.maxCount {
		g.lock.Lock()
		if g.curCount < g.maxCount {
			g.waitGroup.Add(1)
			g.curCount++
			go g.worker()
		}
		g.lock.Unlock()
	}
	g.jobs <- param
}

func (g *GPool) Clear() {
	g.lock.Lock()
	for g.curCount > 0 {
		g.curCount--
		g.jobs <- nil
	}
	g.lock.Unlock()
}

func (g *GPool) Close() {
	close(g.jobs)
	g.waitGroup.Wait()
}

func (g *GPool) worker() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%#v\n", p)
		}
		g.waitGroup.Done()
	}()
	for j := range g.jobs {
		if j == nil {
			break
		}
		g.fun(j)
	}
}
