package docker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thought-machine/falco-probes/pkg/docker"
)

func TestMustClient(t *testing.T) {
	cli := docker.MustClient()

	assert.NotNil(t, cli)
}
