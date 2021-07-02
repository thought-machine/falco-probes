package amazonlinux2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/amazonlinux2"
)

// TODO: optimize these tests, at the moment:
//  - we're running yum install 3 times (and if we carry on like this it will be 6+ times, which makes the test take a v. long time)

func TestGetKernelConfiguration(t *testing.T) {
	cli := docker.MustClient()

	// TODO: this will likely fail in the future when this package is removed from their repositories;
	// 		 we should use a dynamic name and assert it to the best we can.
	kp, err := amazonlinux2.NewKernelPackage(cli, "4.14.200-155.322.amzn2")
	require.NoError(t, err)

	kernelConfigVol := kp.GetKernelConfiguration()

	out, err := cli.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/busybox:latest",
			Entrypoint: []string{"ls"},
			Cmd:        []string{"-l", "/lib/modules/"},
			Volumes: map[operatingsystem.Volume]string{
				kernelConfigVol: "/lib/modules/",
			},
		},
	)
	require.NoError(t, err)
	assert.Contains(t, out, "4.14.200-155.322.amzn2")
}

func TestGetKernelSources(t *testing.T) {
	cli := docker.MustClient()

	// TODO: this will likely fail in the future when this package is removed from their repositories;
	// 		 we should use a dynamic name and assert it to the best we can.
	kp, err := amazonlinux2.NewKernelPackage(cli, "4.14.200-155.322.amzn2")
	require.NoError(t, err)

	kernelSourcesVol := kp.GetKernelSources()

	out, err := cli.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/busybox:latest",
			Entrypoint: []string{"ls"},
			Cmd:        []string{"-l", "/usr/src/kernels/"},
			Volumes: map[operatingsystem.Volume]string{
				kernelSourcesVol: "/usr/src/",
			},
		},
	)
	require.NoError(t, err)
	assert.Contains(t, out, "4.14.200-155.322.amzn2")
}
