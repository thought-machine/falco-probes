package ghreleases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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

	releasesMu sync.Mutex
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
func (ghr *GHReleases) PublishProbe(driverVersion string, probePath string) error {
	probeFileName := filepath.Base(probePath)
	release, err := ghr.ensureReleaseForDriverVersion(driverVersion)
	if err != nil {
		return err
	}

	probeFile, err := os.Open(probePath)
	if err != nil {
		return err
	}
	defer probeFile.Close()

	ctx := context.Background()
	asset, _, err := ghr.ghClient.Repositories.UploadReleaseAsset(ctx, ghr.owner, ghr.repo, release.GetID(), &github.UploadOptions{
		Name: probeFileName,
	}, probeFile)
	if err != nil {
		// TODO: return err when we are using IsAlreadyMirrored()
		log.Warn().
			Str("driver_version", driverVersion).
			Str("probe_file_name", probeFileName).
			Str("path", probePath).
			Err(err).
			Msg("could not upload probe")
		return nil
	}

	log.Info().
		Str("download_url", *asset.BrowserDownloadURL).
		Str("driver_version", driverVersion).
		Str("probe_file_name", probeFileName).
		Str("path", probePath).
		Msg("uploaded probe")

	return nil
}

// IsAlreadyMirrored implmements repository.Repository.IsAlreadyMirrored for GitHub Releases.
func (ghr *GHReleases) IsAlreadyMirrored(driverVersion string, probeName string) (bool, error) {
	// use a mutex to ensure that this function is only called once at a time as it is not thread safe.
	// i.e. a release that doesn't exist may result in multiple goroutines trying to create it at once.
	ghr.releasesMu.Lock()
	defer ghr.releasesMu.Unlock()
	ctx := context.Background()

	// Retrieve the releases
	releases, _, err := ghr.ghClient.Repositories.ListReleases(ctx, ghr.owner, ghr.repo, &github.ListOptions{PerPage: 100}) //TODO: if >100 releases we need to figure out pagination
	if err != nil {
		return false, fmt.Errorf("could not list releases: %w", err)
	}
	for _, release := range releases {
		// Check if release exists for this driverVersion
		if *release.Name == driverVersion {
			log.Info().Str("release", *release.Name).Msg("Matching release found, now checking it's assets: ")
			// Retrieve the matching releases assets
			assets, _, err := ghr.ghClient.Repositories.ListReleaseAssets(ctx, ghr.owner, ghr.repo, *release.ID, &github.ListOptions{PerPage: 100}) //TODO: if >100 assets we need to figure out pagination
			if err != nil {
				return false, fmt.Errorf("could not list release's assets: %w", err)
			}
			for _, asset := range assets {
				// Check if asset matches probeName
				if *asset.Name == probeName {
					log.Info().Str("using", *asset.BrowserDownloadURL).Msg("Probe is uploaded and available")
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("asset not found that matches %s/%s", driverVersion, probeName)
}

func newGHClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func (ghr *GHReleases) ensureReleaseForDriverVersion(driverVersion string) (*github.RepositoryRelease, error) {
	// use a mutex to ensure that this function is only called once at a time as it is not thread safe.
	// i.e. a release that doesn't exist may result in multiple goroutines trying to create it at once.
	ghr.releasesMu.Lock()
	defer ghr.releasesMu.Unlock()

	ctx := context.Background()
	releases, _, err := ghr.ghClient.Repositories.ListReleases(ctx, ghr.owner, ghr.repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list releases: %w", err)
	}
	for _, release := range releases {
		if *release.Name == driverVersion {
			return release, nil
		}
	}

	// release does not exist, create it
	// truncate the driverVersion for the release tag as "branch or tag names consisting of 40 hex characters are not allowed."
	tagName := driverVersion[:8]
	release, _, err := ghr.ghClient.Repositories.CreateRelease(ctx, ghr.owner, ghr.repo, &github.RepositoryRelease{
		Name:    &driverVersion,
		TagName: &tagName,
	})

	return release, err
}
