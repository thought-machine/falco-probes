package docker_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestGetDirectoryFromVolume(t *testing.T) {
	cli := docker.MustClient()

	vol := cli.MustCreateVolume()

	err := cli.WriteFileToVolume(vol, "/tmp/", "/tmp/a", "bar\nbaz")
	require.NoError(t, err)
	err = cli.WriteFileToVolume(vol, "/tmp/", "/tmp/b", "far\nfaz")
	require.NoError(t, err)

	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = cli.GetDirectoryFromVolume(vol, "/tmp/", "/tmp/", tmpDir)
	assert.NoError(t, err)

	files := []string{}
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})
	require.NoError(t, err)

	expectedFiles := []string{
		tmpDir,
		filepath.Join(tmpDir, "tmp"),
		filepath.Join(tmpDir, "tmp/a"),
		filepath.Join(tmpDir, "tmp/b"),
	}

	assert.Equal(t, expectedFiles, files)
}
