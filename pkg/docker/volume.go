package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

const (
	// BusyBoxImage is the image to use for BusyBox, this is used as the smallest Docker image that we can create a container from.
	// We cannot use scratch unless we build our own blank image from scratch, which adds more complexity.
	BusyBoxImage = "docker.io/library/busybox:1.33.1"
)

// MustCreateVolume creates a new docker volume and returns its name, fatally logging any errors.
func (c *Client) MustCreateVolume() operatingsystem.Volume {
	ctx := context.Background()
	vol, err := c.upstream.VolumeCreate(ctx, volume.VolumeCreateBody{})
	if err != nil {
		log.Fatal().Err(err).Msg("could not create docker volume")
	}

	return operatingsystem.Volume(vol.Name)
}

// MustRemoveVolumes removes the given volumes, fatally logging any errors.
func (c *Client) MustRemoveVolumes(volumeNames ...operatingsystem.Volume) {
	ctx := context.Background()
	errs := []error{}
	for _, volumeName := range volumeNames {
		err := c.upstream.VolumeRemove(ctx, string(volumeName), false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		log.Fatal().Errs("errors", errs).Msg("could not remove docker volume(s)")
	}
}

// WriteFileToVolume writes the given contents to the given path with the given volume and its mount point.
// TODO: split into smaller functions
func (c *Client) WriteFileToVolume(volume operatingsystem.Volume, volumeMnt string, path string, contents string) error {
	ctx := context.Background()

	if err := c.EnsureImage(BusyBoxImage); err != nil {
		return err
	}

	volumes := map[operatingsystem.Volume]string{
		volume: volumeMnt,
	}

	resp, err := c.upstream.ContainerCreate(ctx, &container.Config{
		Image:   BusyBoxImage,
		Volumes: getContainerConfigVolumesFromOpts(volumes),
		Tty:     false,
	}, getHostConfigFromOpts(volumes), nil, nil, "")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	header := &tar.Header{
		Name:     filepath.Base(path),
		Mode:     0o777,
		Size:     int64(len(contents)),
		Typeflag: tar.TypeReg,
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		panic(err)
	}
	_, err = tarWriter.Write([]byte(contents))
	if err != nil {
		panic(err)
	}
	err = tarWriter.Close()
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	if err := c.upstream.CopyToContainer(ctx, resp.ID, filepath.Dir(path), reader, types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("could not copy to container: %w", err)
	}

	return c.upstream.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
}

// GetFileFromVolume returns a reader for the contents of a file in the given volume and path.
// TODO: split into smaller functions
func (c *Client) GetFileFromVolume(volume operatingsystem.Volume, volumeMnt string, path string) (io.Reader, error) {
	ctx := context.Background()

	if err := c.EnsureImage(BusyBoxImage); err != nil {
		return nil, err
	}

	volumes := map[operatingsystem.Volume]string{
		volume: volumeMnt,
	}

	resp, err := c.upstream.ContainerCreate(ctx, &container.Config{
		Image:   BusyBoxImage,
		Volumes: getContainerConfigVolumesFromOpts(volumes),
		Tty:     false,
	}, getHostConfigFromOpts(volumes), nil, nil, "")
	if err != nil {
		return nil, err
	}

	reader, _, err := c.upstream.CopyFromContainer(ctx, resp.ID, path)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(reader)
	outBytes := &bytes.Buffer{}
	foundFiles := []string{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}
		foundFiles = append(foundFiles, hdr.Name)
		if _, err := io.Copy(outBytes, tr); err != nil {
			log.Fatal().Err(err)
		}
	}

	if len(foundFiles) != 1 {
		return nil, fmt.Errorf("found more than 1 or no files (%d)", len(foundFiles))
	}

	if err := c.upstream.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		return nil, fmt.Errorf("could not remove container: %w", err)
	}

	return outBytes, nil
}
