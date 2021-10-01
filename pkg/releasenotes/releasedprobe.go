package releasenotes

import (
	"sort"
	"strconv"
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
}

// ReleasedProbeFromMarkdownRow converts a row from a markdown table into a ReleasedProbe
func ReleasedProbeFromMarkdownRow(s string) ReleasedProbe {
	if len(s) < 3 || s[0] != '|' || s[len(s)-1] != '|' {
		return ReleasedProbe{}
	}

	s = s[1 : len(s)-1] // remove preceding and trailing |
	split := strings.Split(s, "|")
	if len(split) != 2 || split[0] == strings.Repeat("-", len(split[0])) || split[0] == " Kernel Package " {
		return ReleasedProbe{}
	}

	return ReleasedProbe{KernelPackage: split[0], Probe: split[1]}
}

// KernelPackageFromProbeName extracts the kernel package from a full probe name
func KernelPackageFromProbeName(probe string) string {
	var kernelPkg string
	switch {
	case strings.Contains(probe, ".amzn2."):
		if probeSplit := strings.Split(probe, "_"); len(probeSplit) > 3 && probeSplit[2] != "" {
			kernelPkg = probeSplit[2] // ... remove the 'falco_amazonlinux2_' prefix
		}

		kernelPkg = kernelPkg[:strings.Index(kernelPkg, ".amzn2.")+6] // ... trim everything from ".amzn2." (inclusive)
	}

	return kernelPkg
}

// ReleasedProbes is a sortable slice of probes, ordered by the most recent kernel package
type ReleasedProbes []ReleasedProbe

var _ sort.Interface = (ReleasedProbes)(nil)

func (rp ReleasedProbes) Len() int      { return len(rp) }
func (rp ReleasedProbes) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }

func (rp ReleasedProbes) Less(i, j int) bool {
	// Split the kernel package down into it's semver elements (assuming . or - as delimiters)
	iSplit := strings.FieldsFunc(rp[i].KernelPackage, splitKernelProbeElements)
	jSplit := strings.FieldsFunc(rp[j].KernelPackage, splitKernelProbeElements)

	for ii := range iSplit {
		if len(jSplit)-1 < ii {
			return false
		}

		if iSplit[ii] == jSplit[ii] {
			continue
		}

		// Convert the semver elements into ints so we can do a numeric '<' comparison (so 2 < 11 == true)
		iInt, errI := strconv.Atoi(iSplit[ii])
		jInt, errJ := strconv.Atoi(jSplit[ii])
		if errI != nil || errJ != nil {
			return iSplit[ii] < jSplit[ii]
		}

		return iInt < jInt
	}

	return rp[i].KernelPackage < rp[j].KernelPackage // fallback case, we shouldn't hit this though
}

func splitKernelProbeElements(r rune) bool { return r == '.' || r == '-' }
