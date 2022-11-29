package buildid

import (
	"fmt"
	"net/http"
	"time"
)

const (
	maxConcurrent = 100
	timeout       = 5 * time.Second
	urlTemplate   = "https://storage.googleapis.com/cos-tools/%s/kernel_commit"
)

// ValidatorInterface is the interface we can override to mock validation.
type ValidatorInterface interface {
	FilterInvalid([]string) ([]string, error)
}

// Validator implements ValidatorInterface
type Validator struct {
	ValidatorInterface

	Client HTTPClient
}

// HTTPClient is an interface we can use for a mock HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ValidatorResult holds information on a build ID's validity and any error raised during validation.
type ValidatorResult struct {
	buildID string
	valid   bool
	err     error
}

// FilterInvalid takes a list of build IDs and returns only ones which are valid releases.
func (v Validator) FilterInvalid(buildIDsIn []string) ([]string, error) {
	if v.Client == nil {
		v.Client = &http.Client{Timeout: timeout}
	}

	buildIDsOut := make([]string, 0)

	// Use a buffered semaphore channel to limit the number of concurrent connections.
	sem := make(chan bool, maxConcurrent)
	results := make(chan ValidatorResult)

	for _, buildID := range buildIDsIn {
		go v.validate(buildID, sem, results)
	}
	for range buildIDsIn {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("could not validate build id %s: %w", result.buildID, result.err)
		}
		if result.valid {
			buildIDsOut = append(buildIDsOut, result.buildID)
		}
	}
	return buildIDsOut, nil
}

func (v Validator) validate(buildID string, sem chan bool, results chan<- ValidatorResult) {
	req, err := http.NewRequest(http.MethodHead, fmt.Sprintf(urlTemplate, buildID), nil)
	if err != nil {
		results <- ValidatorResult{buildID: buildID, valid: false, err: err}
		return
	}
	// Block until there is a free connection.
	sem <- true
	// Run the request.
	res, err := v.Client.Do(req)
	// Release the buffer to allow the next connection.
	<-sem
	if err != nil {
		results <- ValidatorResult{buildID: buildID, valid: false, err: err}
		return
	}
	if res.StatusCode > 299 {
		results <- ValidatorResult{buildID: buildID, valid: false}
		return
	}
	results <- ValidatorResult{buildID: buildID, valid: true}
}
