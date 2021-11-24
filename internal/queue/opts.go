package queue

import "time"

// Opts represents the available options to a Queue.
type Opts struct {
	Buffer        uint64        `long:"buffer" description:"The amount of tasks to buffer in the queue."`
	SentinelTimer time.Duration `long:"sentinel_timer" description:"The amount of time for the sentinel to wait between polls." default:"500ms"`
}
