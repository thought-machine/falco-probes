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
