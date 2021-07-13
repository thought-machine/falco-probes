package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
)

// BuildOpts abstracts the Docker API from the end-user for building Docker images.
type BuildOpts struct {
	Dockerfile string
	BuildArgs  map[string]*string
	Tags       []string
}

// Build builds a docker image with the given options.
func (c *Client) Build(opts *BuildOpts) error {
	ctx := context.Background()

	dockerCtx, err := dockerfileOnlyContext(opts.Dockerfile)
	if err != nil {
		return fmt.Errorf("could not create docker build context: %w", err)
	}

	out, err := c.upstream.ImageBuild(ctx, dockerCtx, types.ImageBuildOptions{
		PullParent: true,
		Remove:     true,
		Dockerfile: "Dockerfile",
		Tags:       opts.Tags,
		BuildArgs:  opts.BuildArgs,
	})
	if err != nil {
		return fmt.Errorf("could not build docker image: %w", err)
	}
	defer out.Body.Close()
	handleBuildOrPullOutput(out.Body, newBuildDebugLogger(*opts))

	return nil
}

// StrPtr returns a pointer to a string.
func StrPtr(str string) *string {
	return &str
}

// dockerfileOnlyContext returns a Docker build context with the given dockerfile contents as a 'Dockerfile' within a build context (tar archive).
func dockerfileOnlyContext(dockerfile string) (io.Reader, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	defer tarWriter.Close()

	header := &tar.Header{
		Name:     "Dockerfile",
		Mode:     0o777,
		Size:     int64(len(dockerfile)),
		Typeflag: tar.TypeReg,
	}
	err := tarWriter.WriteHeader(header)
	if err != nil {
		return nil, err
	}
	_, err = tarWriter.Write([]byte(dockerfile))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// buildDebugLogger implements io.Writer for debug logging of docker build output.
type buildDebugLogger struct {
	opts BuildOpts
}

// newBuildDebugLogger returns a new buildDebugLogger for debug logging of docker build output.
func newBuildDebugLogger(opts BuildOpts) *buildDebugLogger {
	return &buildDebugLogger{
		opts: opts,
	}
}

// Write impelements io.Writer.Write for debug logging of docker build output.
func (l *buildDebugLogger) Write(p []byte) (int, error) {
	logLine := strings.TrimSpace(string(p))
	if len(logLine) > 0 {
		log.Debug().
			Strs("tags", l.opts.Tags).
			Msg(logLine)
	}

	return len(p), nil
}
