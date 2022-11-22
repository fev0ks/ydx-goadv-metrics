package internal

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ExitHandler struct {
	Cancel []context.CancelFunc
	Stop   []chan bool
	Close  []io.Closer
}

func ProperExitDefer(exitHandler *ExitHandler) {
	log.Println("Graceful exit handler is activated")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)
	select {
	case s := <-signals:
		log.Printf("Received a signal '%s'\n", s)
		log.Println("Cancel active contexts")
		for _, cancel := range exitHandler.Cancel {
			cancel()
		}
		log.Println("Stop active goroutines")
		for _, toStop := range exitHandler.Stop {
			toStop <- true
		}
		log.Println("Close active resources")
		for _, toClose := range exitHandler.Close {
			err := toClose.Close()
			if err != nil {
				log.Printf("failed to close an resource: %v\n", err)
			}
		}
	}
}
