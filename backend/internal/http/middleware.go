package http

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.HandlerFunc

func MiddleWareChain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next.ServeHTTP
	}
}

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		next.ServeHTTP(w, r)
	}
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("🎯", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}