package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"worker_pool/internal/handlers"
	"worker_pool/internal/infrastructure/kafka"
	"worker_pool/internal/infrastructure/postgre"
	"worker_pool/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgre.Connect(ctx, "postgresql://user:pass@postgres:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store := postgre.NewTaskStore(db)

	producer, err := kafka.NewProducer(
		[]string{"redpanda:9092"},
		"api",
		"jobs")

	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	taskHandler := handlers.NewTaskHandler(store, producer)

	mux := http.NewServeMux()

	metrics.Init()

	mux.Handle("/tasks", metrics.Middleware("/tasks", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetALl(w, r)
		case http.MethodPost:
			taskHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/tasks/", metrics.Middleware("/tasks/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)
			log.Println("path:", r.URL.Path, idStr, err)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			taskHandler.GetByID(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("server close error: %v", err)
		}
	}

	log.Println("server stopped gracefully")
}
