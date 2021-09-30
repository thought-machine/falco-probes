package ghreleases_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v37/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

func getCachingReleasesTestClient(t *testing.T) *ghreleases.CachingGHReleasesClient {
	token := getGitHubAuthToken(t)
	return ghreleases.NewCachingGHReleasesClient(token)
}

func TestListReleases(t *testing.T) {
	cachingReleasesClient := getCachingReleasesTestClient(t)
	require.Equal(t, uint64(0), cachingReleasesClient.GetAPIRequestsCount())

	initialReleases, err := cachingReleasesClient.ListReleases()
	require.NoError(t, err)

	// run ListReleases 100x (normally this would result in 100 API requests without caching).
	for i := 0; i < 100; i++ {
		releases, err := cachingReleasesClient.ListReleases()
		require.NoError(t, err)
		assert.ElementsMatch(t, initialReleases, releases)
	}

	// assert that we only have 1 API request with caching.
	assert.Equal(t, uint64(1), cachingReleasesClient.GetAPIRequestsCount())
}

func TestListReleaseAssets(t *testing.T) {
	cachingReleasesClient := getCachingReleasesTestClient(t)
	require.Equal(t, uint64(0), cachingReleasesClient.GetAPIRequestsCount())

	// populate cache of releases (1 API request).
	releases, err := cachingReleasesClient.ListReleases()
	require.NoError(t, err)
	require.Equal(t, uint64(1), cachingReleasesClient.GetAPIRequestsCount())

	initialReleaseAssets, err := cachingReleasesClient.ListReleaseAssets(releases[0].GetID())
	require.NoError(t, err)

	// run ListReleaseAssets 100x (normally this would result in 100 API requests without caching).
	for i := 0; i < 100; i++ {
		releaseAssets, err := cachingReleasesClient.ListReleaseAssets(releases[0].GetID())
		require.NoError(t, err)
		assert.ElementsMatch(t, initialReleaseAssets, releaseAssets)
	}

	// assert the number of additional API requests made to matches the number of pages of assets returned.
	assert.Equal(t, 1+len(initialReleaseAssets), cachingReleasesClient.GetAPIRequestsCount())
}

func TestCreateRelease(t *testing.T) {
	cachingReleasesClient := getCachingReleasesTestClient(t)
	require.Equal(t, uint64(0), cachingReleasesClient.GetAPIRequestsCount())

	// populate cache of releases (1 API request).
	_, err := cachingReleasesClient.ListReleases()
	require.NoError(t, err)
	require.Equal(t, uint64(1), cachingReleasesClient.GetAPIRequestsCount())

	// run CreateRelease (1 API request).
	newRelease, err := cachingReleasesClient.CreateRelease(&github.RepositoryRelease{
		Name:    github.String(TestReleaseName),
		TagName: github.String(TestReleaseName),
		// use Draft release to prevent noise.
		Draft: github.Bool(true),
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), cachingReleasesClient.GetAPIRequestsCount())

	// Verify that ListReleases contains the newly created Release (should be cached).
	releases, err := cachingReleasesClient.ListReleases()
	require.NoError(t, err)
	assert.Contains(t, releases, newRelease)

	// assert that we only 2 API requests with caching.
	assert.Equal(t, uint64(2), cachingReleasesClient.GetAPIRequestsCount())

	cleanupTestReleases(t)
}

func TestUploadReleaseAsset(t *testing.T) {
	cachingReleasesClient := getCachingReleasesTestClient(t)
	require.Equal(t, uint64(0), cachingReleasesClient.GetAPIRequestsCount())

	dummyProbeFile, err := ioutil.TempFile("", "falco-dummy-probe.o")
	require.NoError(t, err)
	err = ioutil.WriteFile(dummyProbeFile.Name(), []byte("falco-dummy-probe"), 0644)
	require.NoError(t, err)

	// run CreateRelease so that we have a release to push assets to (1 API request).
	newRelease, err := cachingReleasesClient.CreateRelease(&github.RepositoryRelease{
		Name:    github.String(TestReleaseName),
		TagName: github.String(TestReleaseName),
		// use Draft release to prevent noise
		Draft: github.Bool(true),
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(1), cachingReleasesClient.GetAPIRequestsCount())

	// run UploadReleaseAsset (1 API request).
	newAsset, err := cachingReleasesClient.UploadReleaseAsset(newRelease.GetID(), &github.UploadOptions{
		Name: filepath.Base(dummyProbeFile.Name()),
	}, dummyProbeFile)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), cachingReleasesClient.GetAPIRequestsCount())

	// Verify that ListReleaseAssets has the newly uploaded asset in (should be cached).
	assets, err := cachingReleasesClient.ListReleaseAssets(newRelease.GetID())
	assert.NoError(t, err)
	assert.Contains(t, assets, newAsset)

	// assert that we only 2 API requests with caching.
	assert.Equal(t, uint64(2), cachingReleasesClient.GetAPIRequestsCount())

	cleanupTestReleases(t)
}
