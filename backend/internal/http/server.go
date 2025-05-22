package http

import (
	"context"
	"log"
	"net/http"
	"time"
	"uptime/internal/handler"
)

func StartServer(ctx context.Context, addr string) {
	router := http.NewServeMux()

	router.HandleFunc("GET /monitor", handler.GetAllMonitors)
	router.HandleFunc("GET /monitor/{monitorId}", handler.GetSingleMonitor)
	router.HandleFunc("POST /monitor", handler.CreateMonitor)
	router.HandleFunc("PUT /monitor", handler.UpdateMonitor)
	router.HandleFunc("DELETE /monitor/{monitorId}", handler.DeleteMonitor)

	middleWareChain := MiddleWareChain(
		RequestLoggerMiddleware,
		AuthMiddleware,
	)

	server := &http.Server{
		Addr:           addr,
		Handler:        middleWareChain(router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		<-ctx.Done()
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("🚨 HTTP server Shutdown Failed: %v", err)
		}
		log.Println("HTTP server gracefully stopped")
	}()

	log.Printf("✅ - API is running --> Server listening on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("🚨 HTTP server error: %s", err)
	}
}
