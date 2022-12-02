package buildid_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos/buildid"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos/mock"
)

func TestFilterInvalid(t *testing.T) {
	// See https://cos.googlesource.com/cos/manifest-snapshots/+log/refs/heads/release-R101/
	testBuildIDs := []string{"17162.40.35", "17162.40.34", "17162.40.25", "17162.40.24", "17154.0.0"}
	// See https://cloud.google.com/container-optimized-os/docs/release-notes/m101
	expectedBuildIDs := []string{"17162.40.34", "17162.40.25"}

	// Create the validator.
	validator := buildid.Validator{}

	// Mock the HTTPClient so we don't have to make calls to Google's COS release bucket when testing.
	validator.Client = &mock.HTTPClient{}
	mock.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		re := regexp.MustCompile(`\d+\.\d+\.\d+`)
		buildID := re.FindString(req.URL.Path)
		// Return 200 if in expectedBuildIDs or 404 otherwise.
		statusCode := 404
		for _, expectedBuildID := range expectedBuildIDs {
			if expectedBuildID == buildID {
				statusCode = 200
				break
			}
		}
		return &http.Response{StatusCode: statusCode}, nil
	}

	actualBuildIDs, err := validator.FilterInvalid(testBuildIDs)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedBuildIDs, actualBuildIDs)
}
