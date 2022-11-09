package falcodriverbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
)

func TestBuildImage(t *testing.T) {
	var tests = []struct {
		FalcoVersion     string
		ExpectedImageFQN string
	}{
		{"0.33.0", "docker.io/thoughtmachine/falco-driver-builder:0.33.0"},
		{"0.29.1", "docker.io/thoughtmachine/falco-driver-builder:0.29.1"},
		{"0.29.0", "docker.io/thoughtmachine/falco-driver-builder:0.29.0"},
		{"0.28.1", "docker.io/thoughtmachine/falco-driver-builder:0.28.1"},
		{"0.28.0", "docker.io/thoughtmachine/falco-driver-builder:0.28.0"},
		{"0.27.0", "docker.io/thoughtmachine/falco-driver-builder:0.27.0"},
		{"0.26.2", "docker.io/thoughtmachine/falco-driver-builder:0.26.2"},
		{"0.26.1", "docker.io/thoughtmachine/falco-driver-builder:0.26.1"},
		{"0.26.0", "docker.io/thoughtmachine/falco-driver-builder:0.26.0"},
		{"0.25.0", "docker.io/thoughtmachine/falco-driver-builder:0.25.0"},
		{"0.24.0", "docker.io/thoughtmachine/falco-driver-builder:0.24.0"},
	}

	cli := docker.MustClient()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.FalcoVersion, func(t *testing.T) {
			t.Parallel()
			falcoDriverBuilderImg, err := falcodriverbuilder.BuildImage(cli, tt.FalcoVersion)
			assert.NoError(t, err)
			assert.Equal(t, tt.ExpectedImageFQN, falcoDriverBuilderImg)
		})
	}
}

func TestGetProbePathFromBuildOutput(t *testing.T) {
	var tests = []struct {
		Description       string
		BuildOutput       string
		ExpectedProbeName string
		ExpectedErr       error
	}{
		{"Probe name is found",
			`* Trying to compile the eBPF probe (falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o)
* eBPF probe located in /root/.falco/falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o
* Success: eBPF probe symlinked to /root/.falco/falco-bpf.o`, "/root/.falco/falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o", nil},
		{"Probe name is not found",
			`* Trying to compile the eBPF probe (falco_amazonlinux2_4.14.232-176.381.amzn2.x86_64_1.o)
make[1]: *** /lib/modules/4.14.232-176.381.amzn2.x86_64/build: No such file or directory.  Stop.
make: *** [Makefile:18: all] Error 2
mv: cannot stat '/usr/src/falco-5c0b863ddade7a45568c0ac97d037422c9efb750/bpf/probe.o': No such file or directory
Unable to load the falco eBPF probe`, "", falcodriverbuilder.ErrCouldNotFindProbePathInOutput},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Description, func(t *testing.T) {
			t.Parallel()
			actualProbeName, err := falcodriverbuilder.GetProbePathFromBuildOutput(tt.BuildOutput)
			if tt.ExpectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.ExpectedErr)
			}
			assert.Equal(t, tt.ExpectedProbeName, actualProbeName)
		})
	}
}

func TestExtractProbeFromVolume(t *testing.T) {

}

func TestGetDriverVersion(t *testing.T) {
	var tests = []struct {
		FalcoVersion               string
		ExpectedFalcoDriverVersion string
	}{
		{"0.33.0", "3.0.1+driver"},
		{"0.29.1", "17f5df52a7d9ed6bb12d3b1768460def8439936d"},
		{"0.29.0", "17f5df52a7d9ed6bb12d3b1768460def8439936d"},
		{"0.28.1", "5c0b863ddade7a45568c0ac97d037422c9efb750"},
		{"0.28.0", "5c0b863ddade7a45568c0ac97d037422c9efb750"},
		{"0.27.0", "5c0b863ddade7a45568c0ac97d037422c9efb750"},
		{"0.26.2", "2aa88dcf6243982697811df4c1b484bcbe9488a2"},
		{"0.26.1", "2aa88dcf6243982697811df4c1b484bcbe9488a2"},
		{"0.26.0", "2aa88dcf6243982697811df4c1b484bcbe9488a2"},
		{"0.25.0", "ae104eb20ff0198a5dcb0c91cc36c86e7c3f25c7"},
		{"0.24.0", "85c88952b018fdbce2464222c3303229f5bfcfad"},
	}

	cli := docker.MustClient()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.FalcoVersion, func(t *testing.T) {
			t.Parallel()
			falcoDriverBuilderImg, err := falcodriverbuilder.BuildImage(cli, tt.FalcoVersion)
			require.NoError(t, err)

			actualFalcoDriverVersion, err := falcodriverbuilder.GetDriverVersion(cli, falcoDriverBuilderImg)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedFalcoDriverVersion, actualFalcoDriverVersion)
		})
	}
}
