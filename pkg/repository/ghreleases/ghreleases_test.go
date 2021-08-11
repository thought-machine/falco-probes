package ghreleases_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/v37/github"
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
