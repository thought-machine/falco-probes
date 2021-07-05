package amazonlinux2

import (
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// AmazonLinux2 implements operatingsystem.OperatingSystem for the amazonlinux2.
type AmazonLinux2 struct {
	operatingsystem.OperatingSystem

	dockerClient *docker.Client
}

// NewAmazonLinux2 returns a new amazonlinux2 implementation of operatingsystem.OperatingSystem.
func NewAmazonLinux2(dockerClient *docker.Client) *AmazonLinux2 {
	return &AmazonLinux2{
		dockerClient: dockerClient,
	}
}

// GetKernelPackageNames implements operatingsystem.OperatingSystem.GetKernelPackageNames for the amazonlinux2.
func (s *AmazonLinux2) GetKernelPackageNames() ([]string, error) {
	out, err := s.dockerClient.Run(
		&docker.RunOpts{
			Image:      "docker.io/library/amazonlinux:2",
			Entrypoint: []string{"bash"},
			Cmd:        []string{"-c", "yum --showduplicates list kernel-devel | tail -n+3 | awk '{ print $2 }'"},
		},
	)
	if err != nil {
		return []string{}, err
	}

	out = strings.TrimSpace(out)
	packageNames := strings.Split(out, "\n")

	return packageNames, nil
}

// GetKernelPackageByName implements operatingsystem.OperatingSystem.GetKernelPackageByName for the amazonlinux2.
func (s *AmazonLinux2) GetKernelPackageByName(name string) (*operatingsystem.KernelPackage, error) {
	return NewKernelPackage(s.dockerClient, name)
}
