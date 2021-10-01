package releasenotes_test

import (
	"context"
	"testing"

	"github.com/thought-machine/falco-probes/pkg/releasenotes"

	"github.com/google/go-github/v37/github"
	"github.com/stretchr/testify/assert"
)

func TestSetReleaseNotes(t *testing.T) {
	releases := createStubReleases(
		[]string{
			"falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o",
			"falco_amazonlinux2_4.14.26-54.32.amzn2.x86_64_1.o",
			"falco_amazonlinux2_4.14.33-59.34.amzn2.x86_64_1.o", // out of order
		},
		[]string{
			"falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o",
			"falco_amazonlinux2_4.14.26-54.32.amzn2.x86_64_1.o",
		},
		[]string{
			"falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o",
		},
	)

	sre := stubReleaseEditor{releases: releases}
	err := releasenotes.SetReleaseNotes(context.Background(), releases, &sre)
	assert.NoError(t, err)

	sre.AssertReleaseBody(t, 0, `
# Probes
| Kernel Package | Probe |
|----------------|-------|
|4.14.238-182.422.amzn2|falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o|
|4.14.33-59.34.amzn2|falco_amazonlinux2_4.14.33-59.34.amzn2.x86_64_1.o|
|4.14.26-54.32.amzn2|falco_amazonlinux2_4.14.26-54.32.amzn2.x86_64_1.o|
`)

	sre.AssertReleaseBody(t, 1, `
# Probes
| Kernel Package | Probe |
|----------------|-------|
|4.14.238-182.422.amzn2|falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o|
|4.14.26-54.32.amzn2|falco_amazonlinux2_4.14.26-54.32.amzn2.x86_64_1.o|
`)

	sre.AssertReleaseBody(t, 2, `
# Probes
| Kernel Package | Probe |
|----------------|-------|
|4.14.238-182.422.amzn2|falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o|
`)
}

func createStubReleases(probeByRelease ...[]string) []*github.RepositoryRelease {
	releases := make([]*github.RepositoryRelease, len(probeByRelease))
	for i, r := range probeByRelease {
		i64 := int64(i)
		releases[i] = &github.RepositoryRelease{ID: &i64, Assets: make([]*github.ReleaseAsset, len(r))}

		for ii, probeName := range r {
			p := probeName
			releases[i].Assets[ii] = &github.ReleaseAsset{Name: &p}
		}
	}

	return releases
}

func TestFilterPackages(t *testing.T) {
	probes := releasenotes.ReleasedProbes{
		{KernelPackage: "kp1"},
		{KernelPackage: "kp1"},
		{KernelPackage: "kp1"},

		{KernelPackage: "kp2"},
		{KernelPackage: "kp2"},

		{KernelPackage: "kp3"},
		{KernelPackage: "kp3"},
		{KernelPackage: "kp3"},

		{KernelPackage: "kp4"},
	}

	kpns := []string{"kp1", "kp2", "kp3", "kp4", "kp5"}

	results := probes.ListKernelPackagesToCompile(kpns, 3)

	assert.Equal(t, []string{"kp2", "kp4", "kp5"}, results)
}

func TestParseProbesFromReleaseNotes(t *testing.T) {
	releaseNotes := `
This is some cruft at the start of the release
# Probes
| Kernel Package | Probe |
|----------------|-------|
|4.14.243-185.433.amzn2|falco_amazonlinux2_4.14.243-185.433.amzn2.x86_64_1.o|
|4.14.241-184.433.amzn2|falco_amazonlinux2_4.14.241-184.433.amzn2.x86_64_1.o|
|4.14.238-182.422.amzn2|falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o|
|4.14.238-182.421.amzn2|falco_amazonlinux2_4.14.238-182.421.amzn2.x86_64_1.o|
not a | probe
also|not|a|probe

this is some cruft at the end of the release`

	release := &github.RepositoryRelease{Body: &releaseNotes}

	exp := releasenotes.ReleasedProbes{
		{
			KernelPackage: "4.14.243-185.433.amzn2",
			Probe:         "falco_amazonlinux2_4.14.243-185.433.amzn2.x86_64_1.o",
		},
		{
			KernelPackage: "4.14.241-184.433.amzn2",
			Probe:         "falco_amazonlinux2_4.14.241-184.433.amzn2.x86_64_1.o",
		},
		{
			KernelPackage: "4.14.238-182.422.amzn2",
			Probe:         "falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o",
		},
		{
			KernelPackage: "4.14.238-182.421.amzn2",
			Probe:         "falco_amazonlinux2_4.14.238-182.421.amzn2.x86_64_1.o",
		},
	}

	result := releasenotes.ParseProbesFromReleaseNotes(release)

	assert.Equal(t, exp, result)
}

type stubReleaseEditor struct {
	releases []*github.RepositoryRelease
}

func (sre stubReleaseEditor) EditReleaseNotesByReleaseID(ctx context.Context, rid int64, body string) error {
	for _, r := range sre.releases {
		if r.GetID() != rid {
			continue
		}

		r.Body = &body
		break
	}

	return nil
}

func (sre stubReleaseEditor) AssertReleaseBody(t *testing.T, rid int64, exp string) {
	for _, r := range sre.releases {
		if r.GetID() != rid {
			continue
		}

		assert.Equal(t, exp, r.GetBody())
		break
	}
}
