package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"worker_pool/pkg/metrics"

	"worker_pool/internal/infrastructure/kafka"
	"worker_pool/internal/infrastructure/postgre"

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

	metrics.Init()

	store := postgre.NewTaskStore(db)
	handler := kafka.NewMessageCreateHandler(store)

	consumer, err := kafka.NewConsumer(
		[]string{"redpanda:9092"},
		"task-create-db-writer",
		"task-create-db-writer",
		[]string{"jobs"},
		handler,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:              ":8082",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("start http server on port 8082")
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("kafka consumer started")

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("metrics shutdown failed: %v", err)
	}
}
