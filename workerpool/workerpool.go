package workerpool

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	// goroutine數量
	WorkerNum int

	// 即時分發，做即時回覆不堆積
	JobChan chan func()

	// 等待queue, 為了不卡住主程序，以及處理請求峰值的時候適當的分配工作
	// 當所有工作都忙碌時會先暫時丟到這個queue，等有空位進行分發
	// 請減少了jobchan分配的效能消耗
	// todo: 先不實做，後續補上
	// waitingQueue deque.Deque[func()]

	workerChan chan func()
	StopSignal chan struct{}
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers.
func NewWorkerPool(workerNum int) *WorkerPool {
	return &WorkerPool{
		WorkerNum:  workerNum,
		JobChan:    make(chan func()),
		workerChan: make(chan func(), workerNum),
		StopSignal: make(chan struct{}),
	}
}

// doWork processes jobs and exits when no more jobs are available.
func (wp *WorkerPool) doWork(workerChan chan func(), wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job := <-workerChan:
			if job != nil {
				job()
			}
		case <-wp.StopSignal:
			fmt.Println("dowork exit 接收到stop signal")
			return
		}
	}
}

// Start initializes the worker pool and manages the worker lifecycle.
func (p *WorkerPool) Start() {
	var workerCount int
	// timeout := time.NewTimer(2 * time.Second)
	// defer timeout.Stop()

	for {
		select {
		case job, ok := <-p.JobChan:
			if !ok {
				// wp.stopAllWorkers(workerCount)
				fmt.Println("JobChan closed. Exiting.")
				return
			}
			select {
			case p.workerChan <- job:
			default:
				if workerCount < p.WorkerNum {
					fmt.Println("Got Job.. do worker")
					p.wg.Add(1)
					go p.doWork(p.workerChan, &p.wg)
					workerCount++
				} else {
					// 處理請求峰值的時候適當的分配工作
					// wating queue
					// 先暫時都卡在jobChan中就好
				}
			}
		case <-p.StopSignal:
			fmt.Println("Stop signal received.1111")
			p.wg.Wait()
			close(p.workerChan)
			return
			// case <-timeout.C:
			// 	fmt.Printf("timeout...workerCount: %d, workerChan: %d\n", workerCount, len(wp.workerChan))
			// 	if workerCount > 0 && len(wp.workerChan) == 0 {
			// 		workerCount--
			// 		wp.workerChan <- nil // Signal to stop a worker
			// 	}
			// 	timeout.Reset(1 * time.Second)
			// case <-wp.StopSignal:
			// 	fmt.Println("Stop signal received.")
			// 	wp.stopAllWorkers(workerCount)
			// 	return
			// }
		}
	}

	// Stop all remaining workers as they become ready.
	// for workerCount > 0 {
	// 	wp.workerChan <- nil
	// 	workerCount--
	// }
	// wp.wg.Wait()

	// timeout.Stop()
}

// stopAllWorkers gracefully stops all workers.
func (wp *WorkerPool) stopAllWorkers(workerCount int) {
	for workerCount > 0 {
		wp.workerChan <- nil
		workerCount--
	}
	wp.wg.Wait()
	close(wp.workerChan)
}

// Submit adds a new job to the worker pool.
func (wp *WorkerPool) Submit(job func()) {
	if job != nil {
		wp.JobChan <- job
	}
}

// SubmitWait submits a job and waits for its completion.
func (wp *WorkerPool) SubmitWait(job func()) {
	if job == nil {
		return
	}
	done := make(chan struct{})
	wp.JobChan <- func() {
		job()
		close(done)
	}
	<-done
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Stop gracefully stops the worker pool after all jobs are completed.
func (wp *WorkerPool) Stop() {
	close(wp.JobChan)

	wp.StopSignal <- struct{}{}
	// close(wp.StopSignal)

	wp.wg.Wait()
}

func (wp *WorkerPool) Status(workerCount int) {
	println("WorkerPool Status:")
	println("WorkerNum:", wp.WorkerNum)
	println("JobChan:", len(wp.JobChan))
	println("workerChan:", len(wp.workerChan))
	println("workerCount:", workerCount)
}

// WaitingQueueSize returns the count of tasks in the waiting queue.
// func (p *WorkerPool) WaitingQueueSize() int {
// 	return int(atomic.LoadInt32(&p.waiting))
// }
