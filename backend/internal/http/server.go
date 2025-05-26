package http

import (
	"context"
	"log"
	"net/http"
	"time"
	"uptime/internal/events"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	Address       string
	Router        *http.ServeMux
	PostgresDB    *pgx.Conn
	kafkaProducer *events.KafkaProducer
}

func (s *Server) routes() {
	s.Router.HandleFunc("GET /monitor", s.handleGetAllMonitors)
	s.Router.HandleFunc("GET /monitor/{monitorId}", s.handleGetSingleMonitor)
	s.Router.HandleFunc("POST /monitor", s.handleCreateMonitor)
	s.Router.HandleFunc("PUT /monitor", s.handleUpdateMonitor)
	s.Router.HandleFunc("DELETE /monitor/{monitorId}", s.handleDeleteMonitor)
}

func NewServer(addr string, db *pgx.Conn, kp *events.KafkaProducer) *Server {
	s := &Server{
		Address:       addr,
		Router:        http.NewServeMux(),
		PostgresDB:    db,
		kafkaProducer: kp,
	}
	s.routes()
	return s
}

func (s *Server) Run(ctx context.Context) {
	middleWareChain := MiddleWareChain(
		RequestLoggerMiddleware,
		AuthMiddleware,
	)

	server := &http.Server{
		Addr:           s.Address,
		Handler:        middleWareChain(s.Router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	handleShutdown(ctx, server)

	log.Printf("✅ - API is running --> Server listening on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("🚨 HTTP server error: %s", err)
	}
}

func handleShutdown(ctx context.Context, server *http.Server) {
	go func() {
		<-ctx.Done()
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("🚨 HTTP server Shutdown Failed: %v", err)
		}
		log.Println("HTTP server gracefully stopped")
	}()
}

func StartServer(ctx context.Context, addr string, db *pgx.Conn, kp *events.KafkaProducer) {
	server := NewServer(addr, db, kp)
	server.Run(ctx)
}
