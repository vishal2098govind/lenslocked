package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vishal2098govind/lenslocked/context"
)

type LoggerMiddleware struct {
}

func (lmw *LoggerMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user := context.User(ctx)
		userPrefix := ""
		if user != nil {
			userPrefix = fmt.Sprintf("user:%v", user.Email)
		}
		timestamp := time.Now().UTC().Format("Mon Jan 2 15:04:05 MST 2006")
		timestampStr := fmt.Sprintf("%v", timestamp)

		prefix := fmt.Sprintf("%s | %s | %s | %s | %s | ", r.RemoteAddr, r.Method, r.RequestURI, userPrefix, timestampStr)
		logger := log.New(os.Stdout, prefix, log.Lshortfile)
		ctx = context.WithLogger(ctx, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (lmw *LoggerMiddleware) IPLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("IP: %v\n", ip)
		next.ServeHTTP(w, r)
	})
}
