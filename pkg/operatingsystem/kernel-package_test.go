package operatingsystem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

func TestKernelPackageProbeName(t *testing.T) {
	var tests = []struct {
		kernelPackage     *operatingsystem.KernelPackage
		expectedProbeName string
	}{
		{
			&operatingsystem.KernelPackage{
				OperatingSystem: "ubuntu",
				KernelRelease:   "4.15.0-147-generic",
				KernelVersion:   "#151-Ubuntu SMP Fri Jun 18 19:21:19 UTC 2021",
			},
			"falco_ubuntu_4.15.0-147-generic_151",
		},
		{
			&operatingsystem.KernelPackage{
				OperatingSystem: "amazonlinux2",
				KernelRelease:   "4.14.143-118.123.amzn2.x86_64",
				KernelVersion:   "#1 SMP Thu Sep 12 16:54:23 UTC 2019",
			},
			"falco_amazonlinux2_4.14.143-118.123.amzn2.x86_64_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedProbeName, func(t *testing.T) {
			actualProbeName := tt.kernelPackage.ProbeName()
			assert.Equal(t, tt.expectedProbeName, actualProbeName)
		})
	}
}

func TestKernelPackageFromProbeName(t *testing.T) {
	var tests = []struct {
		probeName        string
		expKernelPackage string
	}{
		{
			probeName:        "falco_amazonlinux2_4.14.101-91.76.amzn2.x86_64_1.o",
			expKernelPackage: "4.14.101-91.76.amzn2",
		},
		{
			probeName:        "falco_lotsofversioning_1.2.3.4.5.6.7.8.9-10.11.amzn2.x86_64_1.o",
			expKernelPackage: "1.2.3.4.5.6.7.8.9-10.11.amzn2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expKernelPackage, func(t *testing.T) {
			result := operatingsystem.KernelPackageFromProbeName(tt.probeName)
			assert.Equal(t, tt.expKernelPackage, result)
		})
	}
}
