package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrErrorsInParameters = errors.New("errors in input parameters")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 || n <= 0 {
		return ErrErrorsInParameters
	}
	var errCount int32
	ch := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for myTask := range ch {
				myErr := myTask()
				if myErr != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, myTask := range tasks {
		ch <- myTask
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
	}
	close(ch)

	wg.Wait()
	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
