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
|4.14.238-182.422|falco_amazonlinux2_4.14.238-182.422.amzn2.x86_64_1.o|
|4.14.33-59.34|falco_amazonlinux2_4.14.33-59.34.amzn2.x86_64_1.o|
|4.14.26-54.32|falco_amazonlinux2_4.14.26-54.32.amzn2.x86_64_1.o|
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
