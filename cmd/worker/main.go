package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"worker_pool/internal/WorkerPool"
	"worker_pool/internal/infrastructure/postgre"
)

func main() {
	intakeCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStartup()

	db, err := postgre.Connect(startupCtx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgre.NewTaskStore(db)
	handler := WorkerPool.NewJobHandler(store)

	workersCount := 10
	batchSize := 10

	jobs := make(chan WorkerPool.JobTask, batchSize)

	var wg sync.WaitGroup

	for w := 1; w <= workersCount; w++ {
		wg.Add(1)
		go handler.Worker(&wg, w, jobs)
	}

	handler.Producer(intakeCtx, jobs, batchSize)

	wg.Wait()

	log.Println("graceful shutdown complete")
}
