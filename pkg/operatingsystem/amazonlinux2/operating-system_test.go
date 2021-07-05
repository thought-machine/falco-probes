package amazonlinux2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/amazonlinux2"
)

func TestGetKernelPackageNames(t *testing.T) {
	cli := docker.MustClient()
	os := amazonlinux2.NewAmazonLinux2(cli)

	res, err := os.GetKernelPackageNames()

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	// TODO: assert that package list looks like kernel package names (e.g. \d.\d.\d-\d.\d.amzn2)
}

func TestGetKernelPackageByName(t *testing.T) {
	cli := docker.MustClient()
	os := amazonlinux2.NewAmazonLinux2(cli)

	// TODO: this will likely fail in the future when this package is removed from their repositories;
	// 		 we should use a dynamic name and assert it to the best we can.
	res, err := os.GetKernelPackageByName("4.14.200-155.322.amzn2")
	assert.NoError(t, err)

	assert.Equal(t, "4.14.200-155.322.amzn2.x86_64", res.KernelRelease)
	assert.Equal(t, "#1 SMP Thu Oct 15 20:11:12 UTC 2020", res.KernelVersion)
	assert.Equal(t, "x86_64", res.KernelMachine)
	assert.Contains(t, res.OSRelease, "NAME=\"Amazon Linux\"")

	assert.NotEmpty(t, res.KernelConfiguration)
	assert.NotEmpty(t, res.KernelSources)
}
