package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// RunOpts abstracts the Docker API away from the end-user.
type RunOpts struct {
	Image      string
	Entrypoint []string
	Cmd        []string
	// Volumes are volumes to mount to the container in the format <volume name>:<mount path>.
	Volumes    map[operatingsystem.Volume]string
	WorkingDir string
}

// Run runs the given docker image, returning its output as a string.
// TODO: break down into smaller functions
func (c *Client) Run(opts *RunOpts) (string, error) {
	ctx := context.Background()

	reader, err := c.upstream.ImagePull(ctx, opts.Image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	io.Copy(os.Stdout, reader)

	if err := reader.Close(); err != nil {
		return "", err
	}

	resp, err := c.upstream.ContainerCreate(ctx, &container.Config{
		Image:      opts.Image,
		Entrypoint: opts.Entrypoint,
		Cmd:        opts.Cmd,
		Volumes:    getContainerConfigVolumesFromOpts(opts),
		Tty:        false,
		WorkingDir: opts.WorkingDir,
	}, getHostConfigFromOpts(opts), nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := c.upstream.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := c.upstream.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := c.upstream.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	containerLogs := HandleContainerLogs(out)

	if err := out.Close(); err != nil {
		return "", err
	}

	if err := c.upstream.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
		return "", err
	}

	return containerLogs, nil
}

func getContainerConfigVolumesFromOpts(opts *RunOpts) map[string]struct{} {
	mountPointSet := map[string]struct{}{}
	for _, mountPoint := range opts.Volumes {
		mountPointSet[mountPoint] = struct{}{}
	}

	return mountPointSet
}

func getHostConfigFromOpts(opts *RunOpts) *container.HostConfig {
	if len(opts.Volumes) < 1 {
		return nil
	}

	binds := []string{}
	for volumeName, mountPoint := range opts.Volumes {
		binds = append(binds, string(volumeName)+":"+mountPoint)
	}

	return &container.HostConfig{
		Binds: binds,
	}
}
