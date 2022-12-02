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
	ToCancel []context.CancelFunc
	ToStop   []chan struct{}
	ToClose  []io.Closer
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
	log.Println("ToCancel active contexts")
	for _, cancel := range eh.ToCancel {
		cancel()
	}
	log.Println("ToStop active goroutines")
	for _, toStop := range eh.ToStop {
		toStop <- struct{}{}
	}
	log.Println("ToClose active resources")
	for _, toClose := range eh.ToClose {
		fmt.Println("kek1")
		err := toClose.Close()
		if err != nil {
			log.Printf("failed to close an resource: %v\n", err)
		}
	}
}
