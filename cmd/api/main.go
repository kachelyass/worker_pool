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
	"worker_pool/internal/app/handlers"
	"worker_pool/internal/infrastructure/kafka"
	"worker_pool/internal/infrastructure/postgre"
	"worker_pool/pkg/metrics/httpmetrics"
	"worker_pool/pkg/metrics/postgresmetrics"

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
		"jobs",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	taskHandler := handlers.NewTaskHandler(store, producer)

	httpmetrics.Init()
	postgresmetrics.Init()

	router := httpmetrics.NewRouter()

	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetALl(w, r)
		case http.MethodPost:
			taskHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			taskHandler.GetByID(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.RawHandle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("server close error: %v", err)
		}
	}
}
