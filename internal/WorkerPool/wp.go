package WorkerPool

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"worker_pool/internal/infrastructure/postgre"
)

type JobTask struct {
	ID          int
	Description string
	Status      string
}

type JobHandler struct {
	store *postgre.TaskStore
}

func NewJobHandler(store *postgre.TaskStore) *JobHandler {
	return &JobHandler{store: store}
}

func waitOrDone(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (h *JobHandler) Producer(ctx context.Context, jobs chan<- JobTask, batchSize int) {
	defer close(jobs)

	for {
		select {
		case <-ctx.Done():
			log.Printf("producer stopped: %v", ctx.Err())
			return
		default:
		}

		ids, err := h.store.ClaimNextIDs(ctx, batchSize)

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
				log.Printf("producer stopped by context: %v", ctx.Err())
				return
			}

			log.Printf("producer: ClaimNextIDs error: %v", err)

			if err := waitOrDone(ctx, 5*time.Second); err != nil {
				log.Printf("producer wait interrupted: %v", err)
				return
			}
			continue
		}

		if len(ids) == 0 {
			log.Printf("producer: no tasks, sleeping...")

			if err := waitOrDone(ctx, 5*time.Second); err != nil {
				log.Printf("producer wait interrupted: %v", err)
				return
			}
			continue
		}

		for _, id := range ids {
			select {
			case <-ctx.Done():
				log.Printf("producer stopped while sending jobs: %v", ctx.Err())
				return
			case jobs <- JobTask{ID: id}:
				log.Printf("job %d sent to channel", id)
			}
		}
	}
}

func (h *JobHandler) Process(ctx context.Context, task JobTask) {
	res, err := h.store.MarkDone(ctx, task.ID)
	if err != nil {
		log.Printf("error processing task %d: %v", task.ID, err)
		return
	}

	log.Printf("marked task %d: %v", task.ID, res)
}

func (h *JobHandler) Worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan JobTask) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopped: %v", id, ctx.Err())
			return

		case j, ok := <-jobs:
			if !ok {
				log.Printf("worker %d: jobs channel closed, exiting", id)
				return
			}

			log.Printf("worker %d processing task %d", id, j.ID)

			if err := waitOrDone(ctx, 3*time.Second); err != nil {
				log.Printf("worker %d interrupted while processing task %d: %v", id, j.ID, err)
				return
			}

			h.Process(ctx, j)
			log.Printf("worker %d finished task %d", id, j.ID)
		}
	}
}

type PoolManager struct {
	mu      sync.Mutex
	handler *JobHandler
	jobs    chan JobTask
	workers map[int]context.CancelFunc
	nextID  int
	wg      sync.WaitGroup
}
