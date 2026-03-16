package woker_pool

import (
	"fmt"
	"net/http"
	"time"
)

func worker(task http.Request, input <-chan int, output chan<- int) { // пример заглушка
	for j := range input {
		fmt.Println("worker", task, "Started", j)
		time.Sleep(time.Second)
		fmt.Println("worker", task, "finished job", j)
	}
}
