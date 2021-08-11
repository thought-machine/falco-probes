package releasenotes

import (
	"sort"
	"strings"
)

// ReleasedProbe represents a single, compiled probe, released as an asset in a Github release
type ReleasedProbe struct {
	KernelPackage string `json:"kernelPackage"`
	Probe         string `json:"probe"`
}

// ToMarkdownRow converts a ReleasedProbe into a the row of a markdown table
func (rp *ReleasedProbe) ToMarkdownRow() string {
	return "|" + rp.KernelPackage + "|" + rp.Probe + "|"
<<<<<<< HEAD
=======
}

// FromMarkdownRow converts a row from a markdown table into a ReleasedProbe
func ReleasedProbeFromMarkdownRow(s string) ReleasedProbe {
	if len(s) < 3 || s[0] != '|' || s[len(s)-1] != '|' {
		return ReleasedProbe{}
	}

	s = s[1 : len(s)-1] // remove preceding and trailing |
	split := strings.Split(s, "|")
	if len(split) != 2 {
		return ReleasedProbe{}
	}

	return ReleasedProbe{KernelPackage: split[1], Probe: split[2]}
>>>>>>> 352e953 (Prefer actual release notes)
}

// KernelPackageFromProbeName extracts the kernel package from a full probe name...
func KernelPackageFromProbeName(probe string) string {
	if !strings.Contains(probe, "amazonlinux2") {
		return "" // other OS aren't currently supported
	}

<<<<<<< HEAD
	kernelPkg := strings.TrimSuffix(probe, ".o") // Remove the file extension ...
	probeSplit := strings.Split(probe, "_")
	kernelPkg = probeSplit[len(probeSplit)-1]
	if len(probeSplit) > 3 && probeSplit[2] != "" {
		kernelPkg = probeSplit[2] // ... remove the 'falco_$os_' prefix and any underscores in the arch suffix
	}

	lastFullStop := strings.LastIndex(kernelPkg, ".")
	if lastFullStop > -1 && len(kernelPkg)-1 > lastFullStop {
		kernelPkg = kernelPkg[:lastFullStop] // ... remove the arch suffix
	}

	return kernelPkg
}

// ReleasedProbes is a sortable slice of probes, ordered by the most recent kernel package
type ReleasedProbes []ReleasedProbe

var _ sort.Interface = (ReleasedProbes)(nil)

func (rp ReleasedProbes) Len() int           { return len(rp) }
func (rp ReleasedProbes) Swap(i, j int)      { rp[i], rp[j] = rp[j], rp[i] }
func (rp ReleasedProbes) Less(i, j int) bool { return rp[i].KernelPackage > rp[j].KernelPackage }
=======
func (rp ReleasedProbes) Len() int           { return len(rp) }
func (rp ReleasedProbes) Swap(i, j int)      { rp[i], rp[j] = rp[j], rp[i] }
func (rp ReleasedProbes) Less(i, j int) bool { return rp[i].KernelPackage < rp[j].KernelPackage }

// KernelPackageFromProbeName extracts the kernel package from a full probe name.
func KernelPackageFromProbeName(probe string) string {
	if !strings.Contains(probe, "amazonlinux2") {
		return "" // Other operating systems not currently supported
	}

	kernelPkg := strings.TrimSuffix(probe, ".o") // Remove the file extension ...
	probeSplit := strings.Split(probe, "_")
	kernelPkg = probeSplit[len(probeSplit)-1]
	if len(probeSplit) > 3 && probeSplit[2] != "" {
		kernelPkg = probeSplit[2] // ... remove the 'falco_$os_' prefix and any underscores in the arch suffix
	}

	lastFullStop := strings.LastIndex(kernelPkg, ".")
	if lastFullStop > -1 && len(kernelPkg)-1 > lastFullStop {
		kernelPkg = kernelPkg[:lastFullStop] // ... remove the arch suffix
	}

	return kernelPkg
}
>>>>>>> 352e953 (Prefer actual release notes)
