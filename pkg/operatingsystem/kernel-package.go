package operatingsystem

import (
	"fmt"
	"regexp"
	"strings"
)

// KernelPackage represents the required inputs for build a Falco Driver for a Kernel Package.
type KernelPackage struct {
	// OperatingSystem is the name of the Operating System that this KernelPackage belongs to.
	OperatingSystem string
	// Name is the name of the KernelPackage from the Operating System's perspective.
	Name string
	// KernelRelease is the value to mock as the output of `uname -r`.
	KernelRelease string
	// KernelVersion is the value to mock as the output of `uname -v`.
	KernelVersion string
	// KernelMachine is the value to mock as the output of `uname -m`.
	KernelMachine string
	// OSRelease is the file contents to use as the mock of `/etc/os-release`.
	OSRelease FileContents
	// KernelConfiguration is the volume to mount as `/host/lib/modules/`.
	KernelConfiguration Volume
	// KernelSources is the volume to mount as `/usr/src/`.
	KernelSources Volume
}

var kernelVersionRe = regexp.MustCompile(`^#(\d+)`)

// ProbeName returns the ProbeName expected by Falco.
// interpreted from: https://github.com/falcosecurity/falco/blob/0.29.1/scripts/falco-driver-loader#L449
func (kp *KernelPackage) ProbeName() string {
	driverName := "falco"
	targetID := kp.OperatingSystem
	kernelRelease := kp.KernelRelease
	// from: $(uname -v | sed 's/#\([[:digit:]]\+\).*/\1/')
	// this sed command is extracting the first set of digits from the KernelVersion after the #. e.g.
	// `#151-Ubuntu SMP Fri Jun 18 19:21:19 UTC 2021` becomes `151`.
	kernelVersion := ""
	matches := kernelVersionRe.FindStringSubmatch(kp.KernelVersion)
	if len(matches) == 2 {
		kernelVersion = matches[1]
	}

	return fmt.Sprintf("%s_%s_%s_%s", driverName, targetID, kernelRelease, kernelVersion)
}

// FileContents represents the contents of a file.
type FileContents string

// Volume represents a reference to a Volume (structured collection of files).
type Volume string

// Validate checks that the probeName contains the given operatingSystem
func (kp *KernelPackage) Validate() error {
	if !strings.Contains(kp.ProbeName(), kp.OperatingSystem) {
		return fmt.Errorf("kernel probe name '%s' does not include operating system '%s'", kp.ProbeName(), kp.OperatingSystem)
	}
	return nil
}
