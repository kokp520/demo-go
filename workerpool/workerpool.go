package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type WorkerPool struct {
	WorkerNum  int
	JobChan    chan Job
	ResultChan chan int
	wg         sync.WaitGroup
}

func NewWorkerPool(workerNum int, chanSize int) *WorkerPool {
	return &WorkerPool{
		WorkerNum:  workerNum,
		JobChan:    make(chan Job, chanSize),
		ResultChan: make(chan int, chanSize),
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for job := range wp.JobChan {
		// do work
		fmt.Printf("worker %d, start job %d\n", id, job.ID)
		time.Sleep(1 * time.Second)
		fmt.Printf("worker %d, finish job %d\n", id, job.ID)
		wp.ResultChan <- job.ID
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.WorkerNum; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) CollectResults() {
	for result := range wp.ResultChan {
		fmt.Printf("collected result %d\n", result)
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.JobChan <- job
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.ResultChan)
}
