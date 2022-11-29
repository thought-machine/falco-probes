package cos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos"
)

func TestGetKernelConfiguration(t *testing.T) {
	cli := docker.MustClient()

	kp, err := cos.NewKernelPackage(cli, "cos-101-17162-40-34")
	require.NoError(t, err)

	out, err := cli.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/busybox:latest",
			Entrypoint: []string{"ls"},
			Cmd:        []string{"-l", "/lib/modules/"},
			Volumes: map[operatingsystem.Volume]string{
				kp.KernelConfiguration: "/lib/modules/",
			},
		},
	)

	require.NoError(t, err)
	assert.Contains(t, out, "5.15.65")
}

// Tests that /usr/src/kernels is present but empty as Falco doesn't require the Google COS sources.
func TestGetKernelSources(t *testing.T) {
	cli := docker.MustClient()

	kp, err := cos.NewKernelPackage(cli, "cos-101-17162-40-34")
	require.NoError(t, err)

	out, err := cli.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/busybox:latest",
			Entrypoint: []string{"ls"},
			Cmd:        []string{"-1A", "/usr/src/kernels/"},
			Volumes: map[operatingsystem.Volume]string{
				kp.KernelSources: "/usr/src/",
			},
		},
	)
	require.NoError(t, err)
	assert.Empty(t, out)
}
