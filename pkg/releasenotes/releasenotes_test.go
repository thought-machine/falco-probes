package releasenotes_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/thought-machine/falco-probes/pkg/releasenotes"

	"github.com/stretchr/testify/assert"
)

func TestToMarkdownString(t *testing.T) {
	rp := releasenotes.ReleasedProbe{DriverVersion: "driver", Probe: "probe", KernelPackage: "kp"}
	assert.Equal(t, rp.ToMarkdownRow(), "|driver|kp|probe|")
}

func TestReleasedProbesOrdering(t *testing.T) {
	rps := releasenotes.ReleasedProbes{
		{
			DriverVersion: "ae104eb2", // 0.25.0,
			Probe:         "4",
			KernelPackage: "aa",
		},
		{
			DriverVersion: "2aa88dcf", //0.26.0
			Probe:         "5",
			KernelPackage: "zz",
		},
		{
			DriverVersion: "85c88952", // 0.24.0
			Probe:         "2",
			KernelPackage: "bb",
		},
		{
			DriverVersion: "85c88952", // 0.24.0
			Probe:         "3",
			KernelPackage: "cc",
		},
		{
			DriverVersion: "17f5df52", // 0.29.1
			Probe:         "7",
			KernelPackage: "aa",
		},
		{
			DriverVersion: "85c88952", // 0.24.0
			Probe:         "1",
			KernelPackage: "ba",
		},
		{
			DriverVersion: "5c0b863d", //0.28.1
			Probe:         "6",
			KernelPackage: "aa",
		},
		{
			DriverVersion: "85c88952", // 0.24.0
			Probe:         "0",
			KernelPackage: "aa",
		},
	}

	sort.Sort(rps)

	for i := 0; i < len(rps); i++ {
		assert.Equal(t, rps[i].Probe, strconv.Itoa(i))
	}
}
