package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/internal/cmd"
)

func TestRunParallelAndCollectErrors(t *testing.T) {
	var tests = []struct {
		description string
		inFnReturns []error
		inLimit     int
		outErrs     []error
	}{
		{
			"no errors",
			[]error{nil, nil, nil, nil},
			3,
			[]error{},
		},
		{
			"1 error",
			[]error{nil, nil, fmt.Errorf("foo"), nil},
			3,
			[]error{fmt.Errorf("foo")},
		},
		{
			"2 errors",
			[]error{nil, fmt.Errorf("bar"), fmt.Errorf("foo"), nil},
			3,
			[]error{fmt.Errorf("bar"), fmt.Errorf("foo")},
		},
		{
			"all errors",
			[]error{fmt.Errorf("bar"), fmt.Errorf("bar"), fmt.Errorf("foo"), fmt.Errorf("bar")},
			3,
			[]error{fmt.Errorf("bar"), fmt.Errorf("bar"), fmt.Errorf("foo"), fmt.Errorf("bar")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			inFns := []func() error{}
			for _, fnRes := range tt.inFnReturns {
				fnRes := fnRes // t.Run introduces concurrency
				inFns = append(inFns, func() error {
					return fnRes
				})
			}

			resErrs := cmd.RunParallelAndCollectErrors(inFns, tt.inLimit)
			// As we're running fns in parallel, there's no guarantees on the order of which they return
			// so we're checking that the length is what we expect and asserting whether everything we want is in the returned errs.
			assert.Len(t, resErrs, len(tt.outErrs))
			for _, fnRes := range tt.outErrs {
				assert.Contains(t, resErrs, fnRes)
			}
		})
	}
}

func TestRunParallelAndCollectErrorsLoops(t *testing.T) {

	amountOfNumbers := 10
	parallelism := 4

	inFns := []func() error{}
	outChan := make(chan int, amountOfNumbers)

	// create functions which output a number to the channel
	for number := 0; number < amountOfNumbers; number++ {
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		number := number

		inFns = append(inFns, func() error {
			outChan <- number
			return nil
		})
	}

	resErrs := cmd.RunParallelAndCollectErrors(inFns, parallelism)
	for _, err := range resErrs {
		assert.NoError(t, err)
	}
	close(outChan)

	// verify that the channel has the same numbers that we expected to output in
	outNumbers := []int{}
	for i := range outChan {
		outNumbers = append(outNumbers, i)
	}

	for i := 0; i < amountOfNumbers; i++ {
		assert.Contains(t, outNumbers, i)
	}

}
