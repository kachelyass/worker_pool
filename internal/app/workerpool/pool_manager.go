package workerpool

import (
	"context"
	"log"
	"sort"
	"sync"
)

type PoolManager struct {
	mu      sync.Mutex
	handler *JobHandler
	jobs    chan JobTask
	workers map[int]context.CancelFunc
	nextID  int
	wg      sync.WaitGroup
}

func NewPoolManager(handler *JobHandler, queueSize int) *PoolManager {
	return &PoolManager{
		handler: handler,
		jobs:    make(chan JobTask, queueSize),
		workers: make(map[int]context.CancelFunc),
	}
}

func (pm *PoolManager) Jobs() chan JobTask {
	return pm.jobs
}

func (pm *PoolManager) Count() int {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return len(pm.workers)
}

func (pm *PoolManager) AddWorkers(n int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i := 0; i < n; i++ {
		workerID := pm.nextID
		pm.nextID++

		ctx, cancel := context.WithCancel(context.Background())
		pm.workers[workerID] = cancel

		pm.wg.Add(1)
		go pm.handler.Worker(ctx, &pm.wg, workerID, pm.jobs)

		log.Printf("pool: worker %d started", workerID)
	}
}

func (pm *PoolManager) RemoveWorkers(n int) {
	pm.mu.Lock()

	if n > len(pm.workers) {
		n = len(pm.workers)
	}

	ids := make([]int, 0, len(pm.workers))
	for id := range pm.workers {
		ids = append(ids, id)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	toStop := ids[:n]
	cancels := make([]context.CancelFunc, 0, n)

	for _, id := range toStop {
		cancels = append(cancels, pm.workers[id])
		delete(pm.workers, id)
	}

	pm.mu.Unlock()

	for i, cancel := range cancels {
		cancel()
		log.Printf("pool: worker %d stop requested", toStop[i])
	}
}

func (pm *PoolManager) SetWorkers(target int) int {
	if target < 1 {
		target = 1
	}

	current := pm.Count()

	switch {
	case target > current:
		pm.AddWorkers(target - current)
	case target < current:
		pm.RemoveWorkers(current - target)
	}

	return pm.Count()
}

func (pm *PoolManager) Wait() {
	pm.wg.Wait()
}
