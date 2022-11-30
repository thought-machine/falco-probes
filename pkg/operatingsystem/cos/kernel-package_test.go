package cos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos"
)

// Tests that the kernel config is downloaded from the fallback filename: x86_64_defconfig.
// See https://cos.googlesource.com/third_party/kernel/+/cb1d3f28467d33c1165181e8ddd2982311fd4687/arch/x86/configs/
// vs. https://cos.googlesource.com/third_party/kernel/+/c19d150c6bd658510ec786390aec80ad476c7578/arch/x86/configs/
func TestGetKernelConfiguration(t *testing.T) {
	cli := docker.MustClient()

	// TODO: mock the http calls.
	kp, err := cos.NewKernelPackage(cli, "cos-89-16108-403-11")
	require.NoError(t, err)

	out, err := cli.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/busybox:latest",
			Entrypoint: []string{"/bin/sh"},
			Cmd:        []string{"-c", "ls /lib/modules && cat /lib/modules/*/config"},
			Volumes: map[operatingsystem.Volume]string{
				kp.KernelConfiguration: "/lib/modules/",
			},
		},
	)

	require.NoError(t, err)
	assert.Contains(t, out, "5.4.104")
	assert.Contains(t, out, "CONFIG_")
}

// Tests that /usr/src/kernels is present but empty as Falco doesn't require the Google COS sources.
func TestGetKernelSources(t *testing.T) {
	cli := docker.MustClient()

	// TODO: mock the http calls.
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
