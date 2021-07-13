package docker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
)

func TestDockerRunLogs(t *testing.T) {
	cli := docker.MustClient()

	out, err := cli.Run(&docker.RunOpts{
		Image: "docker.io/library/alpine:3.14",
		Cmd:   []string{"cat", "/etc/os-release"},
	})

	require.NoError(t, err)

	assert.Equal(t, `NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.14.0
PRETTY_NAME="Alpine Linux v3.14"
HOME_URL="https://alpinelinux.org/"
BUG_REPORT_URL="https://bugs.alpinelinux.org/"
`, out)
}
