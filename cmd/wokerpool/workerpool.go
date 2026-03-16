package wokerpool

import (
	"fmt"
	"time"
)

func workerPool(numJobs int) {
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	for w := 1; w <= numJobs; w++ {
		go worker(w, jobs, results)
	}

}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second * 5)
		fmt.Println("worker", id, "processing job", j)
		results <- j
	}
}
