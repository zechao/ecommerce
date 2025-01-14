package api

import (
	"fmt"
	"net/http"
	"runtime"
)

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				fmt.Printf("recovering from err %v\n %s", err, buf)
				w.Write([]byte(`{"error":"our server got panic"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
