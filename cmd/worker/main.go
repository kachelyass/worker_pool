package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"worker_pool/internal/app/workerpool"
	"worker_pool/internal/handlers"
	"worker_pool/pkg/metrics"

	"worker_pool/internal/infrastructure/postgre"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	intakeCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStartup()

	db, err := postgre.Connect(startupCtx, "postgresql://user:pass@postgres:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	metrics.Init()

	store := postgre.NewTaskStore(db)
	handler := workerpool.NewJobHandler(store)

	workersCount := 10
	batchSize := 10
	queueSize := 100

	pool := workerpool.NewPoolManager(handler, queueSize)
	server := handlers.NewServer(pool)

	pool.AddWorkers(workersCount)

	go handler.Producer(intakeCtx, pool.Jobs(), batchSize)

	mux := http.NewServeMux()
	mux.HandleFunc("/workers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.GetWorkers(w, r)
		case http.MethodPatch:
			server.SetWorkers(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("http server started on :8081")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-intakeCtx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}

	pool.Wait()

	log.Println("graceful shutdown complete")
}
