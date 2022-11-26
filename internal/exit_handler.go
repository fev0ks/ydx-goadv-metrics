package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ExitHandler struct {
	Cancel []context.CancelFunc
	Stop   []chan struct{}
	Close  []io.Closer
}

func ProperExitDefer(exitHandler *ExitHandler) {
	log.Println("Graceful exit handler is activated")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		s := <-signals
		log.Printf("Received a signal '%s'\n", s)
		exitHandler.shutdown()
	}()
}

func (eh *ExitHandler) shutdown() {
	successfullyFinished := make(chan struct{})
	go func() {
		eh.endHeldObjects()
		successfullyFinished <- struct{}{}
	}()
	select {
	case <-successfullyFinished:
		log.Println("System finished work, graceful shutdown")
		os.Exit(0)
	case <-time.After(1 * time.Minute):
		log.Println("System has not shutdown in time '1m', shutdown with interruption")
		os.Exit(1)
	}
}

func (eh *ExitHandler) endHeldObjects() {
	log.Println("Cancel active contexts")
	for _, cancel := range eh.Cancel {
		cancel()
	}
	log.Println("Stop active goroutines")
	for _, toStop := range eh.Stop {
		toStop <- struct{}{}
	}
	log.Println("Close active resources")
	for _, toClose := range eh.Close {
		fmt.Println("kek1")
		err := toClose.Close()
		if err != nil {
			log.Printf("failed to close an resource: %v\n", err)
		}
	}
}
