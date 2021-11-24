package queue

import (
	"sync/atomic"
	"time"

	"github.com/thought-machine/falco-probes/internal/logging"
)

var log = logging.Logger

// Queue represents a queue that is published to asynchronously.
type Queue struct {
	queueCh chan Task

	queuedTasks uint64
	ackedTasks  uint64
}

// NewQueue returns a new Queue.
func NewQueue(opts *Opts) *Queue {
	q := &Queue{
		queueCh:     make(chan Task, opts.Buffer),
		queuedTasks: 0,
		ackedTasks:  0,
	}

	q.sentinel(opts.SentinelTimer)

	return q
}

// Publish publishes the given task to the Queue. This is a separate channel
// so that it can be done asynchronously.
func (q *Queue) Publish(t Task) error {
	atomic.AddUint64(&q.queuedTasks, 1)

	select {
	case q.queueCh <- t:
		return nil
	default:
		return t.Execute(q)
	}
}

// Consume returns the queue as a channel that can be consumed from.
func (q *Queue) Consume() chan Task {

	return q.queueCh
}

// Ack acknowledges that a task has been processed. If all tasks have been processed for a period of time, the sentinel will close the queue.
func (q *Queue) Ack() {
	atomic.AddUint64(&q.ackedTasks, 1)
}

// Close closes the queue publishing channel.
func (q *Queue) Close() {
	close(q.queueCh)
}

func (q *Queue) sentinel(timer time.Duration) {
	queueEmptyCount := 0

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Info().
				Uint64("acked-tasks", q.ackedTasks).
				Uint64("queued-tasks", q.queuedTasks).
				Uint64("remaining-tasks", q.queuedTasks-q.ackedTasks).
				Msg("sentinel")

			if queueEmptyCount >= 3 {
				return
			}
		}
	}()
	go func() {
		for {
			time.Sleep(timer)

			if q.queuedTasks <= q.ackedTasks {
				queueEmptyCount++
			}

			if queueEmptyCount >= 3 {
				log.Debug().Msg("sentinel closing queue")
				q.Close()
				return
			}
		}
	}()
}
