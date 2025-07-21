package middlewares

import (
	"fmt"
	"net/http"
)

func IPLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("IP: %v\n", ip)
		next.ServeHTTP(w, r)
	})
}
