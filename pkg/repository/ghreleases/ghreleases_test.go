package ghreleases_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/go-github/v37/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
	"golang.org/x/oauth2"
)

const (
	// TestReleaseName represents the name to use for testing releases.
	TestReleaseName = "testtest"
)

func getGitHubAuthToken(t *testing.T) string {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skipf("GITHUB_TOKEN not set")
	}
	return token
}

func getGitHubClient(t *testing.T) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGitHubAuthToken(t)},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func cleanupTestReleases(t *testing.T) {
	ghClient := getGitHubClient(t)
	releases, _, err := ghClient.Repositories.ListReleases(context.Background(), ghreleases.GitHubOrganization, ghreleases.GitHubRepository, &github.ListOptions{})
	require.NoError(t, err)
	for _, release := range releases {
		if *release.Name == TestReleaseName {
			_, err := ghClient.Repositories.DeleteRelease(context.Background(), ghreleases.GitHubOrganization, ghreleases.GitHubRepository, *release.ID)
			require.NoError(t, err)
		}
	}
}

func TestEnsureReleaseForDriverVersion(t *testing.T) {
	ghReleasesClient := ghreleases.MustGHReleases(&ghreleases.Opts{
		Token: getGitHubAuthToken(t),
	})

	// Run 10 workers in parallel
	parallelism := 10
	driverVersion := "85c88952b018fdbce2464222c3303229f5bfcfad"

	releaseCh := make(chan *github.RepositoryRelease, parallelism)
	errCh := make(chan error, parallelism)

	var wg sync.WaitGroup
	startCh := make(chan struct{})
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			// wait for start trigger (`close(startCh)`)
			<-startCh
			defer wg.Done()
			release, err := ghReleasesClient.EnsureReleaseForDriverVersion(driverVersion)
			releaseCh <- release
			errCh <- err
		}()
	}
	// trigger the workers
	close(startCh)
	wg.Wait()
	close(releaseCh)
	close(errCh)

	releases := []*github.RepositoryRelease{}
	for release := range releaseCh {
		releases = append(releases, release)
	}

	for _, release := range releases {
		assert.Equal(t, releases[0], release)
	}

	for err := range errCh {
		assert.NoError(t, err)
	}

	cleanupTestReleases(t)
}
