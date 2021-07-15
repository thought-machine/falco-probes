package amazonlinux2

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// Name represents the name of this operating system
const Name = "amazonlinux2"

// AmazonLinux2 implements operatingsystem.OperatingSystem for the amazonlinux2.
type AmazonLinux2 struct {
	operatingsystem.OperatingSystem

	dockerClient *docker.Client
}

// NewAmazonLinux2 returns a new amazonlinux2 implementation of operatingsystem.OperatingSystem.
func NewAmazonLinux2(dockerClient *docker.Client) operatingsystem.OperatingSystem {
	return &AmazonLinux2{
		dockerClient: dockerClient,
	}
}

// GetName implements operatingsystem.OperatingSystem.GetName for the amazonlinux2.
func (s *AmazonLinux2) GetName() string {
	return Name
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

	packageNames = onlyEBPFCompatiblePackageNames(packageNames)

	return packageNames, nil
}

// GetKernelPackageByName implements operatingsystem.OperatingSystem.GetKernelPackageByName for the amazonlinux2.
func (s *AmazonLinux2) GetKernelPackageByName(name string) (*operatingsystem.KernelPackage, error) {
	return NewKernelPackage(s.dockerClient, name)
}

func onlyEBPFCompatiblePackageNames(packageNames []string) []string {
	ebpfCompatibleNames := []string{}
	re := regexp.MustCompile(`^[0-9]+\.[0-9]+`)
	for _, name := range packageNames {
		// extract kernel <major>.<minor> from name
		majorMinor := re.FindString(name)
		majorMinorParts := strings.Split(majorMinor, ".")
		majorStr := majorMinorParts[0]
		minorStr := majorMinorParts[1]

		major, _ := strconv.Atoi(majorStr)
		minor, _ := strconv.Atoi(minorStr)

		if !(major <= 4 && minor < 14) {
			ebpfCompatibleNames = append(ebpfCompatibleNames, name)
		}
	}

	return ebpfCompatibleNames
}
