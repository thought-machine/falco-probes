package cos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos/mock"
)

func TestGetKernelPackageNames(t *testing.T) {
	cli := docker.MustClient()
	os := cos.NewCos(cli)

	// Mock the validator to avoid making calls to Google's COS release bucket when testing.
	cos.BuildIDValidator = mock.BuildIDValidator{}

	res, err := os.GetKernelPackageNames()

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestGetKernelPackageByName(t *testing.T) {
	cli := docker.MustClient()
	os := cos.NewCos(cli)

	res, err := os.GetKernelPackageByName("cos-101-17162-40-34")
	assert.NoError(t, err)

	assert.Equal(t, "5.15.65", res.KernelRelease)
	assert.Equal(t, "#1 SMP Thu Nov 10 10:13:28 UTC 2022", res.KernelVersion)
	assert.Equal(t, "x86_64", res.KernelMachine)
	assert.Contains(t, res.OSRelease, "NAME=\"Container-Optimized OS\"")

	assert.NotEmpty(t, res.KernelConfiguration)
	assert.NotEmpty(t, res.KernelSources)
}

func TestParseVersion(t *testing.T) {
	expectedVersion := &cos.Version{Milestone: 101, BuildID: "17162.40.34"}

	actualVersion, err := cos.ParseVersion("cos-101-17162-40-34")
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, actualVersion)
}
