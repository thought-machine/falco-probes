package queue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/internal/queue"
)

type testTask struct {
	hasRan bool
}

func (t *testTask) Execute(queue *queue.Queue) error {
	defer queue.Ack()
	t.hasRan = true

	return queue.Publish(&testTask2{hasRan: false})
}

type testTask2 struct {
	hasRan bool
}

func (t *testTask2) Execute(queue *queue.Queue) error {
	defer queue.Ack()
	t.hasRan = true

	return nil
}

func TestQueue(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFn := context.WithCancel(ctx)
	var wg sync.WaitGroup
	q := queue.NewQueue(&queue.Opts{
		Buffer:        10,
		SentinelTimer: time.Millisecond * 50,
	})

	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go queue.Worker(ctx, &wg, q)
	}

	tasks := []*testTask{}
	for i := 0; i <= 1000; i++ {
		task := &testTask{hasRan: false}
		tasks = append(tasks, task)
		q.Publish(task)
	}

	_ = cancelFn // call this when interrupting
	wg.Wait()

	for _, task := range tasks {
		assert.True(t, task.hasRan)
	}
}

func TestQueueRecursion(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFn := context.WithCancel(ctx)
	var wg sync.WaitGroup
	q := queue.NewQueue(&queue.Opts{
		Buffer:        10,
		SentinelTimer: time.Millisecond * 50,
	})

	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go queue.Worker(ctx, &wg, q)
	}

	tasks := []*testTask{}
	for i := 0; i <= 1000; i++ {
		task := &testTask{hasRan: false}
		tasks = append(tasks, task)
		q.Publish(task)
	}

	_ = cancelFn // call this when interrupting
	wg.Wait()

	for _, task := range tasks {
		assert.True(t, task.hasRan)
	}

}
