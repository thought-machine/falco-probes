package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
)

// EnsureImage ensures that the given image exists locally in Docker and can be used for creating containers
// by checking if the given image exists, pulling it if not.
func (c *Client) EnsureImage(image string) error {
	if !c.imageExists(image) {
		ctx := context.Background()
		reader, err := c.upstream.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer reader.Close()

		handleBuildOrPullOutput(reader, newPullDebugLogger(image))
	}

	return nil
}

// imageExists returns whether or not the given docker image exists locally.
func (c *Client) imageExists(image string) bool {
	ctx := context.Background()

	if strings.HasPrefix(image, "docker.io/") {
		image = strings.TrimPrefix(image, "docker.io/")
		if strings.HasPrefix(image, "library/") {
			image = strings.TrimPrefix(image, "library/")
		}
	}

	images, err := c.upstream.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("could not list docker images")
		return false
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == image {
				return true
			}
		}
	}

	return false
}

// pullDebugLogger impelements io.Writer for debug logging of docker pull output.
type pullDebugLogger struct {
	image string
}

// newPullDebugLogger returns a new pullDebugLogger for debug logging of docker pull output.
func newPullDebugLogger(image string) *pullDebugLogger {
	return &pullDebugLogger{
		image: image,
	}
}

// Write implements io.Writer.Write for debug logging of docker pull output.
func (l *pullDebugLogger) Write(p []byte) (int, error) {
	logLine := strings.TrimSpace(string(p))
	if len(logLine) > 0 {
		log.Debug().
			Str("image", l.image).
			Msg(logLine)
	}

	return len(p), nil
}
