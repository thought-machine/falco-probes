package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// Client abstracts the docker client into useful functions.
type Client struct {
	upstream *client.Client
}

// MustClient returns a new docker client, fatally logging any errors.
func MustClient() *Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("could not initialise docker client: %s", err)
	}

	return &Client{upstream: cli}
}

// MustCreateVolume creates a new docker volume and returns its name, fatally logging any errors.
func (c *Client) MustCreateVolume() operatingsystem.Volume {
	ctx := context.Background()
	vol, err := c.upstream.VolumeCreate(ctx, volume.VolumeCreateBody{})
	if err != nil {
		log.Fatalf("could not create docker volume: %s", err)
	}

	return operatingsystem.Volume(vol.Name)
}
