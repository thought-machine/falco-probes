package cmd

import (
	"sync"
)

// RunParallelAndCollectErrors runs the given list of functions in parallel with the given parallel limits and returns the errors.
// this is is similar to https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Wait, but returns all of the encountered errors instead of just one.
func RunParallelAndCollectErrors(fns []func() error, limit int) []error {
	limiter := make(chan struct{}, limit)
	var wg sync.WaitGroup
	// errChan to collect errors from all the goroutines.
	errChan := make(chan error, len(fns))

	for _, fn := range fns {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter <- struct{}{}
			defer func() { <-limiter }()
			errChan <- fn()
		}()
	}

	wg.Wait()
	close(errChan)

	errs := []error{}

	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
