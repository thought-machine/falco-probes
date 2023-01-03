package cos

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	// MilestoneMin is the lowest active milestone. Lower value milestones are not built. See
	// https://cloud.google.com/container-optimized-os/docs/concepts/versioning.
	// Set to 85 (despite deprecation to cover versions that are being upgraded).
	MilestoneMin = 85

	milestonePrefix = "origin/release-R"
	buildIDFormat   = "%d.%d.%d"
)

// See https://cos.googlesource.com/cos/manifest-snapshots/+/refs/tags/17162.40.35
func readHashesToBuildIDTags(repository *git.Repository, url string) (map[plumbing.Hash]*object.Tag, error) {
	hashesToBuildIDTags := make(map[plumbing.Hash]*object.Tag)

	tagObjects, err := repository.TagObjects()
	if err != nil {
		return nil, fmt.Errorf("could not list tags for %s: %w", url, err)
	}

	err = tagObjects.ForEach(func(tagObject *object.Tag) error {
		// Filter tag objects not matching Build ID (major.minor.patch).
		n := 0
		_, err := fmt.Sscanf(tagObject.Name, buildIDFormat, &n, &n, &n)
		if err != nil {
			return nil
		}

		hashesToBuildIDTags[tagObject.Target] = tagObject
		return nil
	})
	if err != nil {
		return nil, err
	}

	return hashesToBuildIDTags, nil
}

// See https://cos.googlesource.com/cos/manifest-snapshots/+/refs/heads/release-R101
func readMilestonesToRefs(repository *git.Repository, url string) (map[int]*plumbing.Reference, error) {
	milestonesToRefs := make(map[int]*plumbing.Reference)

	refs, err := repository.References()
	if err != nil {
		return nil, fmt.Errorf("could not list references for %s: %w", url, err)
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		nameShort := ref.Name().Short()
		if strings.HasPrefix(nameShort, milestonePrefix) {
			milestone := 0
			_, err = fmt.Sscanf(nameShort, milestonePrefix+"%d", &milestone)
			if err != nil {
				return err
			}
			milestonesToRefs[milestone] = ref
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return milestonesToRefs, nil
}

// ReadMilestonesToBuildIDs returns a map of milestone with a list of build IDs for that milestone.
// See https://cos.googlesource.com/cos/manifest-snapshots/+log/refs/heads/release-R101/
func ReadMilestonesToBuildIDs(repository *git.Repository, url string) (map[int][]string, error) {
	milestonesToBuildIDs := make(map[int][]string)

	milestonesToRefs, err := readMilestonesToRefs(repository, url)
	if err != nil {
		return nil, fmt.Errorf("could not list milestones for %s: %w", url, err)
	}

	hashesToBuildIDTags, err := readHashesToBuildIDTags(repository, url)
	if err != nil {
		return nil, fmt.Errorf("could not list build ids for %s: %w", url, err)
	}

	for milestone, ref := range milestonesToRefs {
		if milestone < MilestoneMin {
			continue
		}

		log, err := repository.Log(&git.LogOptions{
			From: ref.Hash(),
		})
		if err != nil {
			return nil, fmt.Errorf("could not list milestones for %s: %w", url, err)
		}

		milestonesToBuildIDs[milestone] = make([]string, 0)

		err = log.ForEach(func(commit *object.Commit) error {
			buildIDTag, ok := hashesToBuildIDTags[commit.Hash]
			if ok {
				milestonesToBuildIDs[milestone] = append(milestonesToBuildIDs[milestone], buildIDTag.Name)
			}
			return nil
		})
	}

	return milestonesToBuildIDs, nil
}
