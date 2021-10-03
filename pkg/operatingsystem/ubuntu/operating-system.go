package ubuntu

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// Name represents the name of this operating system.
const Name = "ubuntu"

// Ubuntu implementes operatingsystem.OperatingSystem for ubuntu.
type Ubuntu struct {
	operatingsystem.OperatingSystem

	dockerClient *docker.Client
}

// NewUbuntu returns a new ubuntu implementation of operatingsystem.OperatingSystem.
func NewUbuntu(dockerClient *docker.Client) operatingsystem.OperatingSystem {
	return &Ubuntu{
		dockerClient: dockerClient,
	}
}

// GetName implements operatingsystem.OperatingSystem.GetName for ubuntu.
func (s *Ubuntu) GetName() string {
	return Name
}

// GetKernelPackageNames implements operatingsystem.OperatingSystem.GetKernelPackageNames for ubuntu.
func (s *Ubuntu) GetKernelPackageNames() ([]string, error) {
	ubuntuDownloaderImage, err := BuildUbuntuDownloader(s.dockerClient)
	if err != nil {
		return nil, fmt.Errorf("could not build ubuntudownloader: %w", err)
	}

	out, err := s.dockerClient.Run(
		&docker.RunOpts{
			Image:      ubuntuDownloaderImage,
			Entrypoint: []string{"bash"},
			Cmd:        []string{"-c", "apt-get update > /dev/null 2>&1 && apt-cache search linux-headers | awk '{ print $1 }' | sort -uV"},
		},
	)

	if err != nil {
		return []string{}, err
	}

	out = strings.TrimSpace(out)
	packageNames := strings.Split(out, "\n")

	return onlyEBPFCompatiblePackageNames(packageNames), nil
}

// GetKernelPackageByName implements operatingsystem.OperatingSystem.GetKernelPackageByName for ubuntu.
func (s *Ubuntu) GetKernelPackageByName(name string) (*operatingsystem.KernelPackage, error) {
	return NewKernelPackage(s.dockerClient, name)
}

func onlyEBPFCompatiblePackageNames(packageNames []string) []string {
	ebpfCompatibleNames := []string{}
	re := regexp.MustCompile(`^linux-headers-([0-9]+\.[0-9]+)`)
	for _, name := range packageNames {
		// extract kernel <major>.<minor> from name
		matches := re.FindStringSubmatch(name)
		if len(matches) < 2 {
			// skip meta packages (as they just pin to the latest)
			continue
		}
		majorMinor := matches[1]
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
