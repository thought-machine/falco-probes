package resolver

import (
	"fmt"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/amazonlinux2"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos"
)

// OperatingSystems represents the available operating systems to use and their constructors.
var OperatingSystems = map[string]func(*docker.Client) operatingsystem.OperatingSystem{
	amazonlinux2.Name: amazonlinux2.NewAmazonLinux2,
	cos.Name:          cos.NewCos,
}

// OperatingSystem resolves the given operatingsystem name to an implementation of operatingsystem.OperatingSystem
func OperatingSystem(dockerClient *docker.Client, operatingSystemName string) (operatingsystem.OperatingSystem, error) {
	if constructor, ok := OperatingSystems[operatingSystemName]; ok {
		return constructor(dockerClient), nil
	}

	return nil, fmt.Errorf("unsupported operating system: %s", operatingSystemName)
}
