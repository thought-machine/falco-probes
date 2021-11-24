package queue

import (
	"context"
	"sync"
)

// Worker can be used as a worker to consume and process tasks from the given Queue.
// TODO: add error collection.
func Worker(ctx context.Context, wg *sync.WaitGroup, q *Queue) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-q.Consume():
			if task == nil {
				return
			}
			if err := task.Execute(q); err != nil {
				log.Error().Err(err).Msg("failed task")
			}
		}
	}
}
