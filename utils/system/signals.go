// +build linux darwin

package system

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	go func() {

		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				log.Println("syscall.SIGUSR1")
			case syscall.SIGUSR2:
				log.Println("syscall.SIGUSR1")
			case syscall.SIGTERM:
				gracefulStop(signals)
			default:
				log.Println("Got unregistered signal:", v)
			}
		}
	}()
}
