package docker_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

func TestWriteFileToVolume(t *testing.T) {
	cli := docker.MustClient()

	vol := cli.MustCreateVolume()

	err := cli.WriteFileToVolume(vol, "/var/", "/var/foo", "bar\nbaz")
	require.NoError(t, err)

	out, err := cli.Run(&docker.RunOpts{
		Image:      "docker.io/library/busybox:1.33.1",
		Entrypoint: []string{"cat"},
		Cmd:        []string{"/var/foo"},
		Volumes: map[operatingsystem.Volume]string{
			vol: "/var/",
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "bar\nbaz\n", out)

}

func TestGetFileFromVolume(t *testing.T) {
	cli := docker.MustClient()

	vol := cli.MustCreateVolume()

	err := cli.WriteFileToVolume(vol, "/var/", "/var/foo", "bar\nbaz")
	require.NoError(t, err)

	fileReader, err := cli.GetFileFromVolume(vol, "/var/", "/var/foo")
	assert.NoError(t, err)
	fileBytes, err := ioutil.ReadAll(fileReader)
	require.NoError(t, err)
	assert.Equal(t, []byte("bar\nbaz"), fileBytes)
}
