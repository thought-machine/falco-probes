package resolver

import (
	"fmt"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/amazonlinux2"
)

// OperatingSystem resolves the given operatingsystem name to an implementation of operatingsystem.OperatingSystem
func OperatingSystem(dockerClient *docker.Client, operatingSystemName string) (operatingsystem.OperatingSystem, error) {
	switch operatingSystemName {
	case "amazonlinux2":
		return amazonlinux2.NewAmazonLinux2(dockerClient), nil
	}

	return nil, fmt.Errorf("unsupported operating system: %s", operatingSystemName)
}
