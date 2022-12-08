package closes

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

type (
	ModuleClose struct {
		Name     string
		Priority int
		Func     func()
	}
	closes []ModuleClose
)

var closeHandler closes

const (
	MQPriority     = 100
	GormPriority   = 500
	RedisPriority  = 500
	AliLogPriority = 2000
)

func (c closes) Len() int           { return len(c) }
func (c closes) Less(i, j int) bool { return c[i].Priority < c[j].Priority }
func (c closes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// AddShutdown 增加程序结束时需要关闭的服务
func AddShutdown(c ...ModuleClose) {
	closeHandler = append(closeHandler, c...)
}

// SignalClose 监听信号阻塞关闭
func SignalClose() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	sig := <-c
	fmt.Printf("Got %s signal. Aborting...\n", sig)
	Close()
}

// Close 按照优先级调用关闭方法
func Close() {
	sort.Sort(closeHandler)
	if len(closeHandler) > 0 {
		for _, f := range closeHandler {
			fmt.Printf("Close %s ...\n", f.Name)
			f.Func()
		}
	}
	os.Exit(0)
}
