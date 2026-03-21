package WorkerPool

import (
	"context"
	"fmt"
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

func (h *JobHandler) Producer(ctx context.Context, jobs chan<- JobTask, limit int) {
	for {
		ids, err := h.store.ClaimNextIDs(ctx, limit)
		if len(ids) == 0 {
			fmt.Printf("i go to sleep....ZZZzzzZZZzzz")
			time.Sleep(time.Second * 5)
			continue
		}
		if err != nil {
			log.Printf("producer: ClaimNextIDs: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		for _, id := range ids {
			fmt.Printf("Job ID go to chan: %d\n", id)
			jobs <- JobTask{ID: id}
		}
	}
}

/*
Вообщем продюсер берет таски со статусом пендинг и кладет их сразу в канал
мне показалось, что если он будет смотреть постоянно в бд с помощью цикла, во-первых,
мы уберем задержку, то есть мы сразу будем брать таски и класть их в job по мере их поступления
но с другой стороны постоянно делать запросы в бд как будто слишком дорого
было бы круто если бы был совет как это лучше сделать на проде, но пока что для меня это выглядит как неплохой вариант
*/

func (h *JobHandler) Process(ctx context.Context, task JobTask) {
	res, err := h.store.MarkDone(ctx, task.ID)
	if err != nil {
		log.Printf("Error processing task %d: %v", task.ID, err)
	}
	log.Printf("Marked task %d: %v", task.ID, res)
}

/*
Это я добавил чтобы хендлить изменение статуса на done, допустим, что здесь допустим есть какая - то крутая логика, кроме изменения
одной строчки
*/

func (h *JobHandler) Worker(ctx context.Context, id int, jobs <-chan JobTask) {
	for j := range jobs {
		fmt.Printf("Worker %d processing task %d\n", id, j.ID)
		time.Sleep(time.Second * 3)
		h.Process(ctx, j)
		fmt.Printf("Worker %d finished processing task %d\n", id, j.ID)
	}
} // воркеры которые будут вызывать Process хз норм ли делать для такой таски такую вложенность

type PoolManager struct {
	mu      sync.Mutex // идея в том, чтобы операции обавления воркеров и удаления не мешали друг другу
	handler *JobHandler
	jobs    chan JobTask
	workers map[int]context.CancelFunc
	nextID  int            // чтобы присуждать следующим воркерам уникальный id
	wg      sync.WaitGroup // чтобы дождаться выключения всех работающих воркеров до добавления новых или удаления старыъ
}
