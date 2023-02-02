package middlewares

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
)

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rContentEncoding := r.Header.Get(rest.ContentEncoding)
		if rContentEncoding == "gzip" {
			if r.Body != nil {
				gz, err := zlib.NewReader(r.Body)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer gz.Close()
				body, err := io.ReadAll(gz)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Printf("Length: %d", len(body))
				r.Body = io.NopCloser(strings.NewReader(string(body)))
			}
		} else if rContentEncoding != "" {
			err := fmt.Errorf("unsupported Content-Encoding: %s", rContentEncoding)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotImplemented)
		}
		next.ServeHTTP(w, r)
	})
}

func TimerTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%v] Request time execution for: %s '%s' \r", time.Since(start), r.Method, r.RequestURI)
	})
}
