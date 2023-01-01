package middlewares

import (
	"net/http"
)

type HashChecker struct {
	HashKey string
}

func (hc HashChecker) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if hc.HashKey == "" {
		//	next.ServeHTTP(w, r)
		//} else {
		//	contentType := r.Header.Get(rest.ContentType)
		//	if contentType == rest.ApplicationJSON && r.Body != nil {
		//		body, _ := io.ReadAll(r.Body)
		//		defer r.Body.Close()
		//		err := json.Unmarshal(body, &metric)
		//		if err != nil {
		//			log.Printf("failed to parse metric request: %v\n", err)
		//			http.Error(writer, err.Error(), http.StatusBadRequest)
		//			return
		//		}
		//		decodeString, err := hex.DecodeString(string(bytes))
		//		if err != nil {
		//			return
		//		}
		//		r.Body = io.NopCloser(strings.NewReader(string(body)))
		//	}
		//}

	})
}
