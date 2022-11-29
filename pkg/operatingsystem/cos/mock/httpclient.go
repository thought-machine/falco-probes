package mock

import "net/http"

// HTTPClient is the mock HTTPClient.
type HTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// GetDoFunc fetches the mock HTTPClient's `Do` func.
var GetDoFunc func(req *http.Request) (*http.Response, error)

// Do is the mock HTTPClient's `Do` func.
func (m *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
