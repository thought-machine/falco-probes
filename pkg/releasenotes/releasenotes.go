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

// ReleasedProbes is a sortable slice of probes, sorting by driver ver/kernel pkg ('newest' release last)
type ReleasedProbes []ReleasedProbe

var _ sort.Interface = (ReleasedProbes)(nil)

func (rp ReleasedProbes) Len() int      { return len(rp) }
func (rp ReleasedProbes) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ReleasedProbes) Less(i, j int) bool {
	if rp[i].DriverVersion == rp[j].DriverVersion {
		return rp[i].KernelPackage < rp[j].KernelPackage
	}

	fi, fj := 999, 999 // assume unmatched versions are newer than those listed (so default to 999).
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

// orderedDriverVersions lists the released Falco driver versions by release
var orderedDriverVersions = []string{
	"85c88952", // ~0.24
	"ae104eb2", // ~0.25
	"2aa88dcf", // ~0.26
	"5c0b863d", // ~0.28
	"17f5df52", // ~0.29
}
