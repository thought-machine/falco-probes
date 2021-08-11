package ghreleases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 352e953 (Prefer actual release notes)

	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/repository"

	"github.com/google/go-github/v37/github"
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

	ghClient *CachingGHReleasesClient
}

// MustGHReleases returns a new GitHub Releases repository, fatally erroring if an error is encountered.
func MustGHReleases(opts *Opts) *GHReleases {
	return &GHReleases{
		ghClient: NewCachingGHReleasesClient(opts.Token),
	}
}

// PublishProbe implements repository.Repository.PublishProbe for GitHub Releases.
func (ghr *GHReleases) PublishProbe(driverVersion string, probePath string) error {
	probeFileName := filepath.Base(probePath)
	release, err := ghr.EnsureReleaseForDriverVersion(driverVersion)
	if err != nil {
		return err
	}

	probeFile, err := os.Open(probePath)
	if err != nil {
		return err
	}
	defer probeFile.Close()

	asset, err := ghr.ghClient.UploadReleaseAsset(release.GetID(), &github.UploadOptions{
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

// IsAlreadyMirrored implements repository.Repository.IsAlreadyMirrored for GitHub Releases.
func (ghr *GHReleases) IsAlreadyMirrored(driverVersion string, probeName string) (bool, error) {
	// Retrieve the releases
	release, err := ghr.getReleaseByName(driverVersion)
	if err != nil {
		return false, fmt.Errorf("could not get release: %w", err)
	}
	asset, err := ghr.getAssetFromReleaseByName(release, probeName)
	if err != nil {
		return false, fmt.Errorf("could not get asset: %w", err)
	}
	// log.Info().Str("using", *asset.BrowserDownloadURL).Msg("Probe is uploaded and available")
	log.Info().Msgf("Found probe, access with: curl -LO \"%s\"", *asset.BrowserDownloadURL)
	return true, nil
}

// GetReleases uses the github API to list all previous releases
func (ghr *GHReleases) GetReleases() ([]*github.RepositoryRelease, error) {
	releases, err := ghr.ghClient.ListReleases()
	if err != nil {
		return nil, fmt.Errorf("could not list releases: %w", err)
	}

	return releases, nil
}

// getAssetFromReleaseByName uses the github API to identify whether the desired probe is an asset of the given release
func (ghr *GHReleases) getAssetFromReleaseByName(release *github.RepositoryRelease, probeName string) (*github.ReleaseAsset, error) {

	assets, err := ghr.ghClient.ListReleaseAssets(release.GetID())
	if err != nil {
		return nil, fmt.Errorf("could not list release's assets: %w", err)
	}

	for _, asset := range assets {
		if asset.GetName() == probeName {
			return asset, nil
		}
	}

	return nil, fmt.Errorf("could not find matching asset for: %s", probeName)
}

// getReleaseByName uses the github API to identify the name of the release for the given driverVersion
func (ghr *GHReleases) getReleaseByName(driverVersion string) (*github.RepositoryRelease, error) {

	releases, err := ghr.ghClient.ListReleases()
	if err != nil {
		return nil, fmt.Errorf("could not list releases: %w", err)
	}

	for _, release := range releases {
		if release.GetName() == driverVersion {
			return release, nil
		}
	}

	return nil, fmt.Errorf("could not find matching release for: %s", driverVersion)
}

func newGHClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// EnsureReleaseForDriverVersion creates a release for the given driver version if one does not already exist.
func (ghr *GHReleases) EnsureReleaseForDriverVersion(driverVersion string) (*github.RepositoryRelease, error) {
	if release, err := ghr.getReleaseByName(driverVersion); err == nil {
		return release, nil
	}

	// release does not exist, create it
	// truncate the driverVersion for the release tag as "branch or tag names consisting of 40 hex characters are not allowed."
	tagName := driverVersion[:8]
	release, err := ghr.ghClient.CreateRelease(&github.RepositoryRelease{
		Name:    github.String(driverVersion),
		TagName: github.String(tagName),
	})

	return release, err
}
