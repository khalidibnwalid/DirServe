package middlewares

import "net/http"

type middleware = func(http.Handler) http.Handler


func ApplyMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
