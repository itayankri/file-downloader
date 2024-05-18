package workmanager

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type CONTEXT_KEY string

const WORKER_ID_KEY CONTEXT_KEY = "worker_id"

type WorkManager struct {
	MaxWorkers int
	waitGroup  sync.WaitGroup
}

func (wm *WorkManager) Run(workers []Worker) {
	semaphore := make(chan struct{}, wm.MaxWorkers)
	for _, w := range workers {
		ctx := context.WithValue(context.Background(), WORKER_ID_KEY, uuid.New().String())
		wm.waitGroup.Add(1)
		semaphore <- struct{}{}
		go func(ctx context.Context, worker Worker) error {
			defer wm.waitGroup.Done()
			err := worker.Work(ctx)
			<-semaphore
			return err
		}(ctx, w)
	}
	wm.waitGroup.Wait()
}

func NewWorkManager(maxWorkers int) WorkManager {
	return WorkManager{
		MaxWorkers: maxWorkers,
	}
}

// ====================================================================================================================================================================================

type Worker interface {
	Work(context.Context) error
}
