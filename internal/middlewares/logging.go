package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s - %s %s %s - %s",
			r.Proto,
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
