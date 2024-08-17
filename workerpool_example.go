package main

import (
	"fmt"

	g "demoadi/workerpool"
)

func main() {
	// goroutine.Worker()
	fmt.Println("Hello, world!")

	numWorkers := 3
	numJobs := 10

	workerPool := g.NewWorkerPool(numWorkers, numJobs)

	// Adding jobs to the Job Queue
	for i := 1; i <= numJobs; i++ {
		workerPool.AddJob(g.Job{ID: i})
	}

	close(workerPool.JobChan)

	workerPool.Start()
	workerPool.Wait()
	workerPool.CollectResults()
}
