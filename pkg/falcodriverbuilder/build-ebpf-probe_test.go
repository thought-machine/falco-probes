package falcodriverbuilder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

// Only test a sample of versions as otherwise Github Actions times out.
func TestBuildEBPFProbe(t *testing.T) {
	var tests = []struct {
		falcoVersion        string
		operatingSystemName string
		kernelPackageName   string
	}{
	    // Amazon Linux 2
		{"0.33.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.31.1", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.29.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.26.2", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.24.0", "amazonlinux2", "4.14.200-155.322.amzn2"},

		// Additional tests for extra amzn2 kernels
		{"0.24.0", "amazonlinux2", "4.14.243-185.433.amzn2"},
		{"0.33.0", "amazonlinux2", "4.14.243-185.433.amzn2"},
		{"0.24.0", "amazonlinux2", "5.4.105-48.177.amzn2"},
		{"0.33.0", "amazonlinux2", "5.4.105-48.177.amzn2"},

		// Google COS
		{"0.33.0", "cos", "cos-101-17162-40-34"},
		{"0.31.1", "cos", "cos-101-17162-40-34"},
		{"0.29.0", "cos", "cos-101-17162-40-34"},
		{"0.26.2", "cos", "cos-101-17162-40-34"},
		{"0.24.0", "cos", "cos-101-17162-40-34"},

		// Additional tests for extra cos versions
		{"0.24.0", "cos", "cos-97-16919-0-3"},
		{"0.33.0", "cos", "cos-97-16919-0-3"},
		{"0.24.0", "cos", "cos-93-16623-0-5"},
		{"0.33.0", "cos", "cos-93-16623-0-5"},
	}

	cli := docker.MustClient()

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s-%s-%s", tt.falcoVersion, tt.operatingSystemName, tt.kernelPackageName), func(t *testing.T) {
			t.Parallel()
			operatingSystem, err := resolver.OperatingSystem(cli, tt.operatingSystemName)
			require.NoError(t, err)
			kernelPackage, err := operatingSystem.GetKernelPackageByName(tt.kernelPackageName)
			require.NoError(t, err)
			_, _, err = falcodriverbuilder.BuildEBPFProbe(cli, tt.falcoVersion, operatingSystem, kernelPackage)
			assert.NoError(t, err)
		})

	}
}
