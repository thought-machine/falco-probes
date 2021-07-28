package ghreleases_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/v37/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
	"golang.org/x/oauth2"
)

func getGitHubAuthToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func getGitHubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func TestPublishProbe(t *testing.T) {
	githubAuthToken := getGitHubAuthToken()
	if githubAuthToken == "" {
		t.SkipNow()
	}

	ghReleases := ghreleases.MustGHReleases(&ghreleases.Opts{
		Token: githubAuthToken,
	})

	err := ghReleases.PublishProbe("test", "")
	assert.NoError(t, err)

	ghClient := getGitHubClient(githubAuthToken)
	releases, _, err := ghClient.Repositories.ListReleases(context.Background(), "thought-machine", "falco-probes", &github.ListOptions{})
	require.NoError(t, err)
	for _, release := range releases {
		if *release.Name == "test" {
			_, err := ghClient.Repositories.DeleteRelease(context.Background(), "thought-machine", "falco-probes", *release.ID)
			require.NoError(t, err)
		}
	}
}
