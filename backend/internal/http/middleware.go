package http

import (
	"log"
	"net/http"
	"time"
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
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("🎯 [%s]\t(%v)\t%s", r.Method, time.Since(start), r.URL.Path)
	}
}
