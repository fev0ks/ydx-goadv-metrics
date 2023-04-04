package middlewares

import (
	"errors"
	"log"
	"net/http"

	ext "github.com/fev0ks/ydx-goadv-metrics/internal/shutdown"
)

type ShutdownBlocker struct {
	exitHandler *ext.ExitHandler
}

func NewShutdownBlocker(exitHandler *ext.ExitHandler) *ShutdownBlocker {
	return &ShutdownBlocker{
		exitHandler: exitHandler,
	}
}

func (sb *ShutdownBlocker) BlockTillFinish(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alias := r.URL.Path
		if sb.exitHandler.IsNewFuncExecutionAllowed() {
			sb.exitHandler.AddFuncInProcessing(alias)
			defer sb.exitHandler.FuncFinished(alias)
			next.ServeHTTP(w, r)
		} else {
			log.Println("System is going to shutdown, new func execution is rejected!")
			http.Error(w, errors.New("server is going to shutdown, new requests are not allowed").Error(), http.StatusServiceUnavailable)
		}
	})
}
