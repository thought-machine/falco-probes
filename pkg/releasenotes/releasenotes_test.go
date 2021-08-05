package releasenotes_test

import (
	"testing"

	"github.com/thought-machine/falco-probes/pkg/releasenotes"

	"github.com/stretchr/testify/assert"
)

func TestToMarkdownString(t *testing.T) {
	rp := releasenotes.ReleasedProbe{Probe: "probe", KernelPackage: "kp"}
	assert.Equal(t, rp.ToMarkdownRow(), "|kp|probe|")
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
			probeName:        "falco_notamazon_1.2.3.4.5.6.7.8.9-10.11.amzn2.x86_64_1.o",
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
