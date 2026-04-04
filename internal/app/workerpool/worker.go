package workerpool

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"worker_pool/internal/infrastructure/postgre"
	"worker_pool/pkg/metrics"
)

type JobTask struct {
	ID          int
	Description string
	Status      string
}

type JobHandler struct {
	store       *postgre.TaskStore
	taskTimeout time.Duration
}

func NewJobHandler(store *postgre.TaskStore) *JobHandler {
	return &JobHandler{
		store:       store,
		taskTimeout: 15 * time.Second,
	}
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
			log.Printf("producer: stop taking new tasks: %v", ctx.Err())
			return
		default:
		}

		ids, err := h.store.ClaimNextIDs(ctx, batchSize)
		if err != nil {
			if ctx.Err() != nil ||
				errors.Is(err, context.Canceled) ||
				errors.Is(err, context.DeadlineExceeded) {
				log.Printf("producer: stopped by context: %v", ctx.Err())
				return
			}

			log.Printf("producer: ClaimNextIDs error: %v", err)

			if err := waitOrDone(ctx, 5*time.Second); err != nil {
				log.Printf("producer: stop during backoff: %v", err)
				return
			}
			continue
		}

		if len(ids) == 0 {
			log.Printf("producer: no tasks, sleeping...")

			if err := waitOrDone(ctx, 5*time.Second); err != nil {
				log.Printf("producer: stop during idle wait: %v", err)
				return
			}
			continue
		}

		for _, id := range ids {
			jobs <- JobTask{ID: id}
			metrics.WorkerPoolJobsEnqueuedTotal.Inc()
			metrics.WorkerPoolQueueSize.Set(float64(len(jobs)))
			log.Printf("producer: job %d sent to channel", id)
		}

		select {
		case <-ctx.Done():
			log.Printf("producer: current batch sent, stop taking next tasks: %v", ctx.Err())
			return
		default:
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

func (h *JobHandler) doWork(ctx context.Context, task JobTask) error {
	return waitOrDone(ctx, 3*time.Second)
}

func (h *JobHandler) Worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan JobTask) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d: stop requested", id)
			return

		case job, ok := <-jobs:
			if !ok {
				log.Printf("worker %d: jobs channel closed, exiting", id)
				return
			}

			metrics.WorkerPoolQueueSize.Set(float64(len(jobs)))
			metrics.WorkerPoolJobsStartedTotal.Inc()

			start := time.Now()

			log.Printf("worker %d processing task %d", id, job.ID)

			taskCtx, cancel := context.WithTimeout(context.Background(), h.taskTimeout)

			if err := h.doWork(taskCtx, job); err != nil {
				metrics.WorkerPoolJobsFailedTotal.Inc()
				metrics.WorkerPoolJobDuration.Observe(time.Since(start).Seconds())

				log.Printf("worker %d: task %d failed during work: %v", id, job.ID, err)
				cancel()
				continue
			}

			h.Process(taskCtx, job)
			cancel()

			metrics.WorkerPoolJobsCompletedTotal.Inc()
			metrics.WorkerPoolJobDuration.Observe(time.Since(start).Seconds())
			metrics.WorkerPoolQueueSize.Set(float64(len(jobs)))

			log.Printf("worker %d finished task %d", id, job.ID)
		}
	}
}
