package amazonlinux2

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

var (
	amazonLinux2Image      = "docker.io/library/amazonlinux:2"
	falcoDriverLoaderImage = "docker.io/falcosecurity/falco-driver-loader:0.33.0"
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

	yumDownloaderImage, err := BuildYumDownloader(dockerClient)
	if err != nil {
		return fmt.Errorf("could not build falco-driver-loader: %w", err)
	}

	script := `
set -euo pipefail
yumdownloader kernel-%s kernel-devel-%s
rpm2cpio kernel-%s.$(uname -m).rpm | cpio --extract --make-directories
rpm2cpio kernel-devel-%s.$(uname -m).rpm | cpio --extract --make-directories
`

	_, err = dockerClient.Run(
		&docker.RunOpts{
			Image:      yumDownloaderImage,
			Entrypoint: []string{"/bin/bash"},
			Cmd:        []string{"-c", fmt.Sprintf(script, kp.Name, kp.Name, kp.Name, kp.Name)},
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
	osReleaseVol := dockerClient.MustCreateVolume()
	_, err := dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"cp"},
			Cmd:        []string{"/etc/os-release", "/host/etc/os-release"},
			Volumes: map[operatingsystem.Volume]string{
				osReleaseVol: "/host/etc/",
			},
		},
	)
	if err != nil {
		return err
	}

	fileReader, err := dockerClient.GetFileFromVolume(osReleaseVol, "/host/etc/", "/host/etc/os-release")
	if err != nil {
		return err
	}

	fileContents, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return err
	}

	kp.OSRelease = operatingsystem.FileContents(fileContents)

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
			Image:      amazonLinux2Image,
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
			Image:      falcoDriverLoaderImage,
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
			Image:      falcoDriverLoaderImage,
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
			Image:      falcoDriverLoaderImage,
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
