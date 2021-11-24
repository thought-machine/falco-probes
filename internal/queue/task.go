package queue

// Task abstracts the implementation of tasks that can be placed on a Queue.
type Task interface {
	Execute(queue *Queue) error
}
