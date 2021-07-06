package amazonlinux2

import (
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// NewKernelPackage returns a new hydrated example implementation operatingsystem.KernelPackage.
func NewKernelPackage(dockerClient *docker.Client, name string) (*operatingsystem.KernelPackage, error) {
	kP := &operatingsystem.KernelPackage{
		OperatingSystem: "amazonlinux2",
		Name:            name,
	}

	if err := addSourcesAndConfiguration(dockerClient, kP); err != nil {
		return nil, err
	}

	if err := addOSRelease(dockerClient, kP); err != nil {
		return nil, err
	}

	if err := addKernelReleaseAndVersionAndMachine(dockerClient, kP); err != nil {
		return nil, err
	}

	return kP, nil
}

func addSourcesAndConfiguration(dockerClient *docker.Client, kp *operatingsystem.KernelPackage) error {
	kp.KernelConfiguration = dockerClient.MustCreateVolume()
	kp.KernelSources = dockerClient.MustCreateVolume()

	// TODO: use http.Get lib to download from repository directly and extract into a docker volume.
	_, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"yum"},
			Cmd:        []string{"-y", "install", "kernel-devel-" + kp.Name, "kernel-" + kp.Name},
			Volumes: map[operatingsystem.Volume]string{
				kp.KernelSources:       "/usr/src/",
				kp.KernelConfiguration: "/lib/modules/",
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func addOSRelease(dockerClient *docker.Client, kp *operatingsystem.KernelPackage) error {
	out, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"cat"},
			Cmd:        []string{"/etc/os-release"},
		},
	)
	if err != nil {
		return err
	}
	kp.OSRelease = operatingsystem.FileContents(out)

	return nil
}

func addKernelReleaseAndVersionAndMachine(dockerClient *docker.Client, kp *operatingsystem.KernelPackage) error {
	kernelSrcPath, err := findKernelSrcPath(dockerClient, kp.KernelSources, kp.Name)
	if err != nil {
		return err
	}

	kernelRelease, err := getKernelRelease(dockerClient, kp.KernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.KernelRelease = kernelRelease

	kernelVersion, err := getKernelVersion(dockerClient, kp.KernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.KernelVersion = kernelVersion

	kernelMachine, err := getKernelMachine(dockerClient, kp.KernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.KernelMachine = kernelMachine

	return nil
}

func findKernelSrcPath(dockerClient *docker.Client, kernelSrcsVol operatingsystem.Volume, name string) (string, error) {
	out, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"find"},
			Cmd:        []string{"/usr/src/", "-name", "*" + name + "*", "-type", "d"},
			Volumes: map[operatingsystem.Volume]string{
				kernelSrcsVol: "/usr/src/",
			},
		},
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func getKernelRelease(dockerClient *docker.Client, kernelSrcsVol operatingsystem.Volume, kernelSrcPath string) (string, error) {
	out, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/falcosecurity/falco-driver-loader:0.28.1",
			Entrypoint: []string{"/bin/bash"},
			Cmd:        []string{"-c", "make kernelrelease | tail -n1"},
			Volumes: map[operatingsystem.Volume]string{
				kernelSrcsVol: "/usr/src/",
			},
			WorkingDir: kernelSrcPath,
		},
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func getKernelVersion(dockerClient *docker.Client, kernelSrcsVol operatingsystem.Volume, kernelSrcPath string) (string, error) {
	out, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/falcosecurity/falco-driver-loader:0.28.1",
			Entrypoint: []string{"/bin/bash"},
			Cmd:        []string{"-c", "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_VERSION.* | cut -f2 -d\\\""},
			Volumes: map[operatingsystem.Volume]string{
				kernelSrcsVol: "/usr/src/",
			},
			WorkingDir: kernelSrcPath,
		},
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func getKernelMachine(dockerClient *docker.Client, kernelSrcsVol operatingsystem.Volume, kernelSrcPath string) (string, error) {
	out, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/falcosecurity/falco-driver-loader:0.28.1",
			Entrypoint: []string{"/bin/bash"},
			Cmd:        []string{"-c", "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_MACHINE.* | cut -f2 -d\\\""},
			Volumes: map[operatingsystem.Volume]string{
				kernelSrcsVol: "/usr/src/",
			},
			WorkingDir: kernelSrcPath,
		},
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
