package shutdown

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	mu = &sync.Mutex{}
)

type ExitHandler struct {
	server            *http.Server
	ToCancel          []context.CancelFunc
	ToStop            []chan struct{}
	ToClose           []io.Closer
	ToExecute         []func() error
	funcsInProcessing sync.WaitGroup
	newFuncAllowed    bool
}

func NewExitHandler() *ExitHandler {
	return &ExitHandler{
		newFuncAllowed:    true,
		funcsInProcessing: sync.WaitGroup{},
	}
}

func (eh *ExitHandler) IsNewFuncExecutionAllowed() bool {
	mu.Lock()
	defer mu.Unlock()
	return eh.newFuncAllowed
}

func (eh *ExitHandler) setNewFuncExecutionAllowed(value bool) {
	mu.Lock()
	defer mu.Unlock()
	eh.newFuncAllowed = value
}

func (eh *ExitHandler) ShutdownServerBeforeExit(server *http.Server) {
	eh.server = server
}

func (eh *ExitHandler) AddFuncInProcessing(alias string) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("'%s' func is started and added to exit handler", alias)
	eh.funcsInProcessing.Add(1)
}

func (eh *ExitHandler) FuncFinished(alias string) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("'%s' func is finished and removed from exit handler", alias)
	eh.funcsInProcessing.Add(-1)
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
		log.Printf("Received a signal '%s'", s)
		exitHandler.setNewFuncExecutionAllowed(false)
		exitHandler.shutdown()
	}()
}

func (eh *ExitHandler) shutdown() {
	successfullyFinished := make(chan struct{})
	go func() {
		eh.waitForShutdownServer()
		eh.waitForFinishFunc()
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

func (eh *ExitHandler) waitForFinishFunc() {
	log.Println("Waiting for functions finish work...")
	eh.funcsInProcessing.Wait()
	log.Println("All functions finished work successfully")
}

func (eh *ExitHandler) waitForShutdownServer() {
	if eh.server != nil {
		log.Println("Waiting for shutdown server...")
		err := eh.server.Shutdown(context.Background())
		log.Println("Server shutdown complete")
		if err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	}
}

func (eh *ExitHandler) endHeldObjects() {
	log.Println("ToExecute final funcs")
	for _, execute := range eh.ToExecute {
		err := execute()
		if err != nil {
			log.Printf("func error: %v", err)
		}
	}
	log.Println("ToCancel active contexts")
	for _, cancel := range eh.ToCancel {
		cancel()
	}
	log.Println("ToStop active goroutines")
	for _, toStop := range eh.ToStop {
		close(toStop)
	}
	log.Println("ToClose active resources")
	for _, toClose := range eh.ToClose {
		err := toClose.Close()
		if err != nil {
			log.Printf("failed to close an resource: %v", err)
		}
	}
	log.Println("Success end final work")
}
