package ghreleases

import (
	"context"
	"fmt"

	"github.com/google/go-github/v37/github"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/repository"
	"golang.org/x/oauth2"
)

var log = logging.Logger

// Opts represents the available options for the GitHub Releases client.
type Opts struct {
	Token string `long:"token" description:"The token to use to authenticate against github" env:"GITHUB_TOKEN" required:"true"`
}

// GHReleases implements repository.Repository against Github Releases.
type GHReleases struct {
	repository.Repository

	ghClient *github.Client
	owner    string
	repo     string
}

// MustGHReleases returns a new GitHub Releases repository, fatally erroring if an error is encountered.
func MustGHReleases(opts *Opts) *GHReleases {

	ghClient := newGHClient(opts.Token)

	return &GHReleases{
		ghClient: ghClient,
		owner:    "thought-machine",
		repo:     "falco-probes",
	}
}

// PublishProbe implmements repository.Repository.PublishProbe for GitHub Releases.
func (ghr *GHReleases) PublishProbe(driverVersion string, probeName string, probePath string) error {
	// TODO: unimplimented
	return fmt.Errorf("unimplemented")
}

// IsAlreadyMirrored implmements repository.Repository.IsAlreadyMirrored for GitHub Releases.
func (ghr *GHReleases) IsAlreadyMirrored(driverVersion string, probeName string) error {
	// TODO: unimplimented
	return fmt.Errorf("unimplemented")
}

func newGHClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
