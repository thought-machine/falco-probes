package docker

import (
	"github.com/docker/docker/client"
	"github.com/thought-machine/falco-probes/internal/logging"
)

var log = logging.Logger

// Client abstracts the docker client into useful functions.
type Client struct {
	upstream *client.Client
}

// MustClient returns a new docker client, fatally logging any errors.
func MustClient() *Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialise docker client")
	}

	return &Client{upstream: cli}
}
