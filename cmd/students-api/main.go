package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaaaxxx/students-api/internal/config"
	"github.com/kaaaxxx/students-api/internal/http/handlers/student"
	"github.com/kaaaxxx/students-api/internal/storage/sqlite"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	//load config
	cfg := config.MustLoad()
	//setup databse

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Db.Close()

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// CORS middleware
	handler := corsMiddleware(router)
	// setup server

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: handler,
	}
	slog.Info("Server started", slog.String("Address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	// Check what these do on internet (os.Interrupt, syscall.SIGINT, syscall.SIGTERM).

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	//Learn more about context and its method.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Server shutdown successfully!")

}
