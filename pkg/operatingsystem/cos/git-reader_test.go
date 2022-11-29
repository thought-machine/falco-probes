package cos_test

import (
	"regexp"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos"
)

func TestReadMilestonesToBuildIDs(t *testing.T) {
	// See https://cos.googlesource.com/cos/manifest-snapshots/+log/refs/heads/release-R101/
	expectedMilestone := 101
	unexpectedMilestone := cos.MilestoneMin - 1

	url := "https://cos.googlesource.com/cos/manifest-snapshots"
	// TODO: Mock the repository.
	repository, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	assert.NoError(t, err)

	milestonesToBuildIDs, err := cos.ReadMilestonesToBuildIDs(repository, url)
	assert.NoError(t, err)

	assert.Contains(t, milestonesToBuildIDs, expectedMilestone)
	assert.NotContains(t, milestonesToBuildIDs, unexpectedMilestone)
	assert.Regexp(t, regexp.MustCompile(`^\d+\.\d+\.\d+$`), milestonesToBuildIDs[expectedMilestone][0])
}
