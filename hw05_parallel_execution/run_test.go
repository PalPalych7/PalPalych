package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		startGoCount := runtime.NumGoroutine()
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")

		finishGoCount := runtime.NumGoroutine()
		require.Equal(t, startGoCount, finishGoCount, "there is no worked gorutin")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		startGoCount := runtime.NumGoroutine()
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")

		finishGoCount := runtime.NumGoroutine()
		require.Equal(t, startGoCount, finishGoCount, "there is no worked gorutin")
	})

	t.Run("tasks count less workersCount", func(t *testing.T) {
		startGoCount := runtime.NumGoroutine()
		tasksCount := 10
		workersCount := 5
		maxErrorsCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
		finishGoCount := runtime.NumGoroutine()
		require.Equal(t, startGoCount, finishGoCount, "there is no worked gorutin")
	})

	t.Run("workersCount 0 or less", func(t *testing.T) {
		tasksCount := 5
		workersCount := 2
		maxErrorsCount := 0
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return nil
			})
		}
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Error(t, err)
	})

	t.Run("tasksCount 0 or less", func(t *testing.T) {
		tasksCount := 3
		workersCount := 0
		maxErrorsCount := 1
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return nil
			})
		}
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Error(t, err)
	})
}
