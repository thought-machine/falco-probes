package releasenotes_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/thought-machine/falco-probes/pkg/releasenotes"

	"github.com/stretchr/testify/assert"
)

func TestToMarkdownString(t *testing.T) {
	rp := releasenotes.ReleasedProbe{Probe: "probe", KernelPackage: "kp"}
	assert.Equal(t, rp.ToMarkdownRow(), "|kp|probe|")
}

func TestReleasedProbeFromMarkdownRow(t *testing.T) {
	rp := releasenotes.ReleasedProbe{Probe: "a", KernelPackage: "hi"}
	mdr := rp.ToMarkdownRow()
	assert.Equal(t, rp, releasenotes.ReleasedProbeFromMarkdownRow(mdr))

	emptyReleaseProbe := releasenotes.ReleasedProbe{}
	assert.Equal(t, emptyReleaseProbe, releasenotes.ReleasedProbeFromMarkdownRow("just a random string"))
	assert.Equal(t, emptyReleaseProbe, releasenotes.ReleasedProbeFromMarkdownRow("|nope|"))
	assert.Equal(t, emptyReleaseProbe, releasenotes.ReleasedProbeFromMarkdownRow("|too|many|rows|in|this|table|"))
}

func TestKernelPackageFromProbeName(t *testing.T) {
	var tests = []struct {
		probeName        string
		expKernelPackage string
	}{
		{
			probeName:        "falco_amazonlinux2_4.14.101-91.76.amzn2.x86_64_1.o",
			expKernelPackage: "4.14.101-91.76.amzn2",
		},
		{
			probeName:        "falco_amazonlinux2_1.2.3.4.5.6.7.8.9-10.11.amzn2.x86_64_1.o",
			expKernelPackage: "1.2.3.4.5.6.7.8.9-10.11.amzn2",
		},
		{
			probeName:        "falco_notamazon_1.2.3.4.5.6.7.8.9-10.11.ubuntu.x86_64_1.o",
			expKernelPackage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expKernelPackage, func(t *testing.T) {
			result := releasenotes.KernelPackageFromProbeName(tt.probeName)
			assert.Equal(t, tt.expKernelPackage, result)
		})
	}
}

func TestSortReleasedProbes(t *testing.T) {
	type testCase struct {
		desc  string
		input releasenotes.ReleasedProbes
	}

	testCases := []testCase{
		{
			"standard sort",
			releasenotes.ReleasedProbes{
				{KernelPackage: "1.bb", Probe: "2"},
				{KernelPackage: "1.baa.ab", Probe: "4"},
				{KernelPackage: "2.aa", Probe: "0"},
				{KernelPackage: "1.ba.ba", Probe: "6"},
				{KernelPackage: "1.baa.ba", Probe: "3"},
				{KernelPackage: "1.cb", Probe: "1"},
				{KernelPackage: "1.ba.ba.a", Probe: "5"},
			},
		},
		{
			"already sorted",
			releasenotes.ReleasedProbes{
				{KernelPackage: "2.aa", Probe: "0"},
				{KernelPackage: "1.cb", Probe: "1"},
				{KernelPackage: "1.bb", Probe: "2"},
				{KernelPackage: "1.baa.ba", Probe: "3"},
				{KernelPackage: "1.baa.ab", Probe: "4"},
				{KernelPackage: "1.ba.ba.a", Probe: "5"},
				{KernelPackage: "1.ba.ba", Probe: "6"},
			},
		},
		{
			"less doesn't panic if comparison is against two elements with different lengths",
			releasenotes.ReleasedProbes{
				{KernelPackage: "1.bb.ab", Probe: "0"},
				{KernelPackage: "1.bb", Probe: "1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			sort.Sort(sort.Reverse(&tc.input))

			for i := 0; i < len(tc.input); i++ {
				assert.Equal(t, tc.input[i].Probe, strconv.Itoa(i))
			}
		})
	}
}
