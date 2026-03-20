package main

import (
	"context"
	"log"
	"time"
	"worker_pool/internal/WorkerPool"
	"worker_pool/internal/infrastructure/postgre"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgre.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgre.NewTaskStore(db)
	handler := WorkerPool.NewJobHandler(store)

	limit := 10

	jobs := make(chan WorkerPool.JobTask, limit)

	for w := 1; w <= limit; w++ {
		go handler.Worker(ctx, w, jobs)
	}

	handler.Producer(ctx, jobs, limit)

}
