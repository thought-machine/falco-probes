package ghreleases

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/go-github/v37/github"
)

const (
	// GitHubOrganization is the GitHub Organization to use.
	GitHubOrganization = "thought-machine"
	// GitHubRepository is the GitHub Repository within the Organization to use.
	GitHubRepository = "falco-probes"
)

var (
	// ErrReleaseNotFound is returned when a requested release could not be found.
	ErrReleaseNotFound = errors.New("release not found")
)

// CachingGHReleasesClient represents a GitHub Releases client which caches the GitHub Releases API.
// in order to reduce API requests.
type CachingGHReleasesClient struct {
	ghClient func() *github.Client

	releasesByID         map[int64]*github.RepositoryRelease
	assetsByID           map[int64]*github.ReleaseAsset
	assetIDsToReleaseIDs map[int64]int64
	dataMu               sync.RWMutex

	apiRequestsCounter uint64

	owner string
	repo  string
}

// NewCachingGHReleasesClient returns a new caching GitHub Releases client.
func NewCachingGHReleasesClient(token string) *CachingGHReleasesClient {
	client := &CachingGHReleasesClient{
		apiRequestsCounter:   0,
		releasesByID:         map[int64]*github.RepositoryRelease{},
		assetsByID:           map[int64]*github.ReleaseAsset{},
		assetIDsToReleaseIDs: map[int64]int64{},
		owner:                GitHubOrganization,
		repo:                 GitHubRepository,
	}
	client.ghClient = func() *github.Client {
		client.apiRequestsCounter++
		return newGHClient(token)
	}
	return client
}

// GetAPIRequestsCount returns the amount of API requests that have been made with this client.
func (c *CachingGHReleasesClient) GetAPIRequestsCount() uint64 {
	return c.apiRequestsCounter
}

// UploadReleaseAsset uploads the given asset to the given release ID.
func (c *CachingGHReleasesClient) UploadReleaseAsset(releaseID int64, opts *github.UploadOptions, file *os.File) (*github.ReleaseAsset, error) {
	release, err := c.GetReleaseByID(releaseID)
	if err != nil {
		return nil, err
	}

	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	ctx := context.Background()
	asset, _, err := c.ghClient().Repositories.UploadReleaseAsset(ctx, c.owner, c.repo, release.GetID(), opts, file)
	if err != nil {
		return nil, err
	}

	c.assetsByID[asset.GetID()] = asset
	c.assetIDsToReleaseIDs[asset.GetID()] = release.GetID()

	return asset, nil
}

// CreateRelease creates a release with the given options.
func (c *CachingGHReleasesClient) CreateRelease(opts *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	ctx := context.Background()
	release, _, err := c.ghClient().Repositories.CreateRelease(ctx, c.owner, c.repo, opts)
	if err != nil {
		return nil, err
	}

	c.releasesByID[release.GetID()] = release

	return release, nil
}

// ListReleaseAssets returns a list of assets for the given GitHub Release ID.
func (c *CachingGHReleasesClient) ListReleaseAssets(releaseID int64) ([]*github.ReleaseAsset, error) {
	release, err := c.GetReleaseByID(releaseID)
	if err != nil {
		return nil, err
	}

	c.dataMu.RLock()

	releaseAssets := []*github.ReleaseAsset{}
	for assetID, associatedReleaseID := range c.assetIDsToReleaseIDs {
		if associatedReleaseID == releaseID {
			releaseAssets = append(releaseAssets, c.assetsByID[assetID])
		}
	}

	if len(releaseAssets) == 0 {
		c.dataMu.RUnlock()
		if err := c.populateUpstreamReleaseAssets(release.GetID()); err != nil {
			return nil, err
		}

		return c.ListReleaseAssets(releaseID)
	}

	c.dataMu.RUnlock()
	return releaseAssets, nil
}

// ListReleases returns the list of Releases.
func (c *CachingGHReleasesClient) ListReleases() ([]*github.RepositoryRelease, error) {
	c.dataMu.RLock()

	resReleases := []*github.RepositoryRelease{}
	for _, release := range c.releasesByID {
		resReleases = append(resReleases, release)
	}

	if len(resReleases) == 0 {
		c.dataMu.RUnlock()
		if err := c.populateUpstreamReleases(); err != nil {
			return nil, err
		}

		return c.ListReleases()
	}

	c.dataMu.RUnlock()
	return resReleases, nil
}

// GetReleaseByID returns a GitHub release by its ID if it exists.
func (c *CachingGHReleasesClient) GetReleaseByID(releaseID int64) (*github.RepositoryRelease, error) {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()

	if release, ok := c.releasesByID[releaseID]; ok {
		return release, nil
	}

	return nil, ErrReleaseNotFound
}

// EditReleaseNotesByReleaseID edits the body/text content of a github release
func (c *CachingGHReleasesClient) EditReleaseNotesByReleaseID(ctx context.Context, releaseID int64, body string) error {
	r := &github.RepositoryRelease{Body: &body}
	if _, _, err := c.ghClient().Repositories.EditRelease(ctx, c.owner, c.repo, releaseID, r); err != nil {
		return err
	}

	c.dataMu.Lock()
	if r := c.releasesByID[releaseID]; r != nil {
		r.Body = &body
	}
	c.dataMu.Unlock()

	return nil
}

// populateUpstreamReleases populates the cache with releases from upstream.
func (c *CachingGHReleasesClient) populateUpstreamReleases() error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	opt := &github.ListOptions{PerPage: 100}
	ctx := context.Background()
	c.releasesByID = map[int64]*github.RepositoryRelease{}
	for {
		releases, releaseResponse, err := c.ghClient().Repositories.ListReleases(ctx, c.owner, c.repo, opt)
		if err != nil {
			return fmt.Errorf("could not list releases: %w", err)
		}
		for _, release := range releases {
			c.releasesByID[release.GetID()] = release
		}
		if releaseResponse.NextPage == 0 {
			break
		}
		opt.Page = releaseResponse.NextPage
	}

	return nil
}

// populateUpstreamReleaseAssets populates the cache with release assets from upstream.
func (c *CachingGHReleasesClient) populateUpstreamReleaseAssets(releaseID int64) error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	opt := &github.ListOptions{PerPage: 100}
	ctx := context.Background()
	for {
		assets, assetResponse, err := c.ghClient().Repositories.ListReleaseAssets(ctx, c.owner, c.repo, releaseID, opt)

		if err != nil {
			return fmt.Errorf("could not list release's assets: %w", err)
		}

		for _, asset := range assets {
			c.assetsByID[asset.GetID()] = asset
			c.assetIDsToReleaseIDs[asset.GetID()] = releaseID
		}

		if assetResponse.NextPage == 0 {
			break
		}
		opt.Page = assetResponse.NextPage
	}

	return nil
}
