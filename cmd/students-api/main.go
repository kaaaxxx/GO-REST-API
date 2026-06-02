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

	"github.com/kaaaxxx/students-api/internal/http/handlers/student"
	"github.com/kaaaxxx/students-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//setup databse 
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// setup server

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("Address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	// Check what these do on internet (os.Interrupt, syscall.SIGINT, syscall.SIGTERM).

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Faield to start server!")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	//Learn more about context and its method.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully!")

}
