package main

import (
	"context"
	"log"
	"sync"
	"time"

	"worker_pool/internal/infrastructure/postgre"
)

func main() {
	ctx := context.Background()

	db, err := postgre.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgre.NewTaskStore(db)

	ids, err := store.GetPendingIDs(ctx)

	const numWorkers = 5
	jobs := make(chan int, len(ids))

	var wg sync.WaitGroup
	wg.Add(len(ids))

	// старт воркеров
	for w := 0; w < numWorkers; w++ {
		go worker(w, jobs, store, ctx, &wg)
	}

	// отправляем задачи
	for _, id := range ids {
		jobs <- id
	}
	close(jobs)

	// ждём, пока все задачи будут обработаны
	wg.Wait()
	log.Println("all pending tasks processed")
}

func worker(workerID int, jobs <-chan int, store *postgre.TaskStore, ctx context.Context, wg *sync.WaitGroup) {
	for taskID := range jobs {
		log.Printf("worker %d starting job %d", workerID, taskID)

		time.Sleep(8 * time.Second)

		_, err := store.UpdateTaskStatus(ctx, taskID)
		if err != nil {
			log.Printf("worker %d error on job %d: %v", workerID, taskID, err)
			wg.Done()
			continue
		}

		log.Printf("worker %d finishing job %d", workerID, taskID)
		wg.Done()
	}
}
