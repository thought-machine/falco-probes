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

func TestBuildEBPFProbe(t *testing.T) {
	var tests = []struct {
		falcoVersion        string
		operatingSystemName string
		kernelPackageName   string
	}{
		{"0.29.1", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.29.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.28.1", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.28.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.27.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.26.2", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.26.1", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.26.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.25.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
		{"0.24.0", "amazonlinux2", "4.14.200-155.322.amzn2"},
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
