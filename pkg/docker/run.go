package docker

import (
	"context"
	"fmt"
	"strings"

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
	Env        map[string]string
}

// Run runs the given docker image, returning its output as a string.
// TODO: break down into smaller functions
func (c *Client) Run(opts *RunOpts) (containerOut string, err error) {
	ctx := context.Background()

	log.Debug().
		Strs("entrypoint", opts.Entrypoint).
		Strs("cmd", opts.Cmd).
		Str("image", opts.Image).
		Msg("docker run")

	if err := c.EnsureImage(opts.Image); err != nil {
		return "", err
	}

	resp, err := c.upstream.ContainerCreate(ctx, &container.Config{
		Image:      opts.Image,
		Entrypoint: opts.Entrypoint,
		Cmd:        opts.Cmd,
		Volumes:    getContainerConfigVolumesFromOpts(opts.Volumes),
		Tty:        true,
		WorkingDir: opts.WorkingDir,
		Env:        envMapToSlice(opts.Env),
	}, getHostConfigFromOpts(opts.Volumes), nil, nil, "")
	if err != nil {
		return "", err
	}

	defer func() {
		if remErr := c.upstream.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); remErr != nil {
			err = remErr
		}
	}()

	if err := c.upstream.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	out, err := c.upstream.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return "", err
	}

	containerOut = handleContainerLogs(out, newRunDebugLogger(*opts))

	if err := out.Close(); err != nil {
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

	containerInspect, err := c.upstream.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", err
	}

	if containerInspect.State.ExitCode != 0 {
		return containerOut, fmt.Errorf("non-zero exit-code (%d) for: %s\n%s",
			containerInspect.State.ExitCode,
			strings.Join(append(opts.Entrypoint, opts.Cmd...), " "),
			containerOut,
		)
	}

	return containerOut, nil
}

func getContainerConfigVolumesFromOpts(volumes map[operatingsystem.Volume]string) map[string]struct{} {
	mountPointSet := map[string]struct{}{}
	for _, mountPoint := range volumes {
		mountPointSet[mountPoint] = struct{}{}
	}

	return mountPointSet
}

func getHostConfigFromOpts(volumes map[operatingsystem.Volume]string) *container.HostConfig {
	if len(volumes) < 1 {
		return nil
	}

	binds := []string{}
	for volumeName, mountPoint := range volumes {
		binds = append(binds, string(volumeName)+":"+mountPoint)
	}

	return &container.HostConfig{
		Binds: binds,
	}
}

func envMapToSlice(envMap map[string]string) []string {
	envSlice := []string{}
	for k, v := range envMap {
		envSlice = append(envSlice, k+"="+v)
	}

	return envSlice
}

type runDebugLogger struct {
	opts RunOpts
}

func newRunDebugLogger(opts RunOpts) *runDebugLogger {
	return &runDebugLogger{
		opts: opts,
	}
}

func (l *runDebugLogger) Write(p []byte) (int, error) {
	logLine := strings.TrimSpace(string(p))
	if len(logLine) > 0 {
		log.Debug().
			Str("image", l.opts.Image).
			Strs("entrypoint", l.opts.Entrypoint).
			Strs("cmd", l.opts.Cmd).
			Msg(logLine)
	}

	return len(p), nil
}
