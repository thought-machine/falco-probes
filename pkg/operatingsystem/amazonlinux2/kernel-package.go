package amazonlinux2

import (
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// KernelPackage implements operatingsystem.KernelPackage for the example.
type KernelPackage struct {
	operatingsystem.KernelPackage
	dockerClient *docker.Client

	name string

	kernelSources       operatingsystem.Volume
	kernelConfiguration operatingsystem.Volume
	osRelease           operatingsystem.FileContents

	kernelRelease string
	kernelVersion string
	kernelMachine string
}

// NewKernelPackage returns a new hydrated example implementation operatingsystem.KernelPackage.
func NewKernelPackage(dockerClient *docker.Client, name string) (*KernelPackage, error) {
	kP := &KernelPackage{
		name:         name,
		dockerClient: dockerClient,
	}

	if err := kP.fetchSourcesAndConfiguration(); err != nil {
		return nil, err
	}

	if err := kP.fetchOSRelease(); err != nil {
		return nil, err
	}

	if err := kP.setKernelReleaseAndVersionAndMachine(); err != nil {
		return nil, err
	}

	return kP, nil
}

// GetKernelRelease implements operatingsystem.KernelPackage.GetKernelRelease for the example.
func (kp *KernelPackage) GetKernelRelease() string {
	return kp.kernelRelease
}

// GetKernelVersion implements operatingsystem.KernelPackage.GetKernelVersion for the example.
func (kp *KernelPackage) GetKernelVersion() string {
	return kp.kernelVersion
}

// GetKernelMachine implements operatingsystem.KernelPackage.GetKernelMachine for the example.
func (kp *KernelPackage) GetKernelMachine() string {
	return kp.kernelMachine
}

// GetOSRelease implements operatingsystem.KernelPackage.GetOSRelease for the example.
func (kp *KernelPackage) GetOSRelease() operatingsystem.FileContents {
	return kp.osRelease
}

// GetKernelConfiguration implements operatingsystem.KernelPackage.GetKernelConfiguration for the example.
func (kp *KernelPackage) GetKernelConfiguration() operatingsystem.Volume {
	return kp.kernelConfiguration
}

// GetKernelSources implements operatingsystem.KernelPackage.GetKernelSources for the example.
func (kp *KernelPackage) GetKernelSources() operatingsystem.Volume {
	return kp.kernelSources
}

func (kp *KernelPackage) fetchSourcesAndConfiguration() error {
	kp.kernelSources = kp.dockerClient.MustCreateVolume()
	kp.kernelConfiguration = kp.dockerClient.MustCreateVolume()

	// TODO: use http.Get lib to download from repository directly and extract into a docker volume.
	_, err := kp.dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"yum"},
			Cmd:        []string{"-y", "install", "kernel-devel-" + kp.name, "kernel-" + kp.name},
			Volumes: map[operatingsystem.Volume]string{
				kp.kernelSources:       "/usr/src/",
				kp.kernelConfiguration: "/lib/modules/",
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (kp *KernelPackage) fetchOSRelease() error {
	out, err := kp.dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"cat"},
			Cmd:        []string{"/etc/os-release"},
		},
	)
	if err != nil {
		return err
	}
	kp.osRelease = operatingsystem.FileContents(out)

	return nil
}

func (kp *KernelPackage) setKernelReleaseAndVersionAndMachine() error {
	kernelSrcPath, err := findKernelSrcPath(kp.dockerClient, kp.kernelSources, kp.name)
	if err != nil {
		return err
	}

	kernelRelease, err := getKernelRelease(kp.dockerClient, kp.kernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.kernelRelease = kernelRelease

	kernelVersion, err := getKernelVersion(kp.dockerClient, kp.kernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.kernelVersion = kernelVersion

	kernelMachine, err := getKernelMachine(kp.dockerClient, kp.kernelSources, kernelSrcPath)
	if err != nil {
		return err
	}
	kp.kernelMachine = kernelMachine

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
