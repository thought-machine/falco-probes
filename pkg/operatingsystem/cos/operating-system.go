package cos

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos/buildid"
)

const (
	// Name represents the name of this operating system.
	Name = "cos"

	// URL versions is the Google COS repo with all the COS milestone and build id information.
	urlVersions = "https://cos.googlesource.com/cos/manifest-snapshots"
)

// BuildIDValidator does what it says on the tin.
var BuildIDValidator buildid.ValidatorInterface

// Cos implements operatingsystem.OperatingSystem for cos.
type Cos struct {
	operatingsystem.OperatingSystem

	dockerClient *docker.Client
}

// Version represents the milestone and build id of a Google COS version which can be found
// in the Image Name (e.g. cos-101-17162-40-34 -> 101, 17162.40.34).
type Version struct {
	Milestone int
	BuildID   string
}

// NewCos returns a new cos implementation of operatingsystem.OperatingSystem.
func NewCos(dockerClient *docker.Client) operatingsystem.OperatingSystem {
	return &Cos{
		dockerClient: dockerClient,
	}
}

// GetName implements operatingsystem.OperatingSystem.GetName for cos.
func (s *Cos) GetName() string {
	return Name
}

// GetKernelPackageNames implements operatingsystem.OperatingSystem.GetKernelPackageNames for cos.
// Note that KernelPackageNames in this context are Google COS Image Names.
func (s *Cos) GetKernelPackageNames() ([]string, error) {
	packageNames := make([]string, 0)

	repository, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: urlVersions,
	})

	milestonesToBuildIDs, err := ReadMilestonesToBuildIDs(repository, urlVersions)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve milestones and build ids: %w", err)
	}

	if BuildIDValidator == nil {
		BuildIDValidator = buildid.Validator{}
	}

	for milestone, candidateBuildIDs := range milestonesToBuildIDs {
		validBuildIDs, err := BuildIDValidator.FilterInvalid(candidateBuildIDs)
		if err != nil {
			return nil, fmt.Errorf("could not filter invalid build ids: %w", err)
		}

		for _, buildID := range validBuildIDs {
			packageNames = append(packageNames, fmt.Sprintf("cos-%d-%s", milestone, strings.ReplaceAll(buildID, ".", "-")))
		}
	}

	return packageNames, nil
}

// GetKernelPackageByName implements operatingsystem.OperatingSystem.GetKernelPackageByName for cos.
func (s *Cos) GetKernelPackageByName(name string) (*operatingsystem.KernelPackage, error) {
	return NewKernelPackage(s.dockerClient, name)
}

// ParseVersion takes an image name (e.g. "cos-101-17162-40-34") and returns the milestone (eg. 101) and build ID (e.g.
// "17162.40.34")
func ParseVersion(name string) (*Version, error) {
	milestone := 0
	buildNumber := ""
	_, err := fmt.Sscanf(name, "cos-%d-%s", &milestone, &buildNumber)
	if err != nil {
		return nil, fmt.Errorf("could not parse version from image name %s: %w", name, err)
	}
	return &Version{Milestone: milestone, BuildID: strings.ReplaceAll(buildNumber, "-", ".")}, nil
}
