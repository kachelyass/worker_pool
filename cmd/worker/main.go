package main

import (
	"context"
	"log"
	"sync"
	"time"

	"worker_pool/internal/WorkerPool"
	"worker_pool/internal/infrastructure/postgre"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := postgre.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
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
		go handler.Worker(ctx, &wg, w, jobs)
	}

	handler.Producer(ctx, jobs, batchSize)

	wg.Wait()

	log.Println("graceful shutdown complete")
}
