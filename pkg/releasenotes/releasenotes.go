package releasenotes

import (
	"sort"
)

// ReleasedProbe represents a single, compiled probe, released as an asset in a Github release
type ReleasedProbe struct {
	DriverVersion string `json:"driverVersion"`
	Probe         string `json:"probe"`
	KernelPackage string `json:"kernelPackage"`
}

// ToMarkdownRow converts a ReleasedProbe into a the row of a markdown table
func (rp *ReleasedProbe) ToMarkdownRow() string {
	return "|" + rp.DriverVersion + "|" + rp.KernelPackage + "|" + rp.Probe + "|"
}

// ReleasedProbes is a sortable slice of probes, ordered by most recent driver, then most recent kernel package
type ReleasedProbes []ReleasedProbe

var _ sort.Interface = (ReleasedProbes)(nil)

func (rp ReleasedProbes) Len() int      { return len(rp) }
func (rp ReleasedProbes) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ReleasedProbes) Less(i, j int) bool {
	// If the Falco versions are the same, sort based on the most recent kernel package...
	if rp[i].DriverVersion == rp[j].DriverVersion {
		return rp[i].KernelPackage > rp[j].KernelPackage
	}

	// Otherwise, we sort by the most recent driver version
	// We assume unmatched versions are newer that those listed (so default to 999).
	fi, fj := 999, 999
	for ii, o := range orderedDriverVersions {
		if o == rp[i].DriverVersion {
			fi = ii
		}
		if o == rp[j].DriverVersion {
			fj = ii
		}
	}

	return fi < fj
}

// This list of driver versions is sorted by the _newest first_.
// Prepend new drivers versions to the top of this list.
// We use the short versions because that is how we name our releases
var orderedDriverVersions = []string{
	"17f5df52",
	"5c0b863d",
	"2aa88dcf",
	"ae104eb2",
	"85c88952",
}
