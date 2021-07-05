package falcodriverbuilder

import (
	// embed is used for including assets via Go 1.16
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

// FalcoDriverBuilderDockerfile contains the Dockerfile contents for build a falco-driver-builder image.
//go:embed falco-driver-builder.Dockerfile
var FalcoDriverBuilderDockerfile string

const (
	// FalcoDriverBuilderRepository is the repository to build the falco-driver-builder image under.
	FalcoDriverBuilderRepository = "docker.io/thoughtmachine/falco-driver-builder"
	// BuiltFalcoProbesDir references the directory where the falco-driver-builder image outputs built probes to.
	BuiltFalcoProbesDir = "/root/.falco/"
)

var (
	// ErrCouldNotFindProbeNameInOutput is returned when a probe could not be found in the output text.
	ErrCouldNotFindProbeNameInOutput = errors.New("could not find build probe name in output")
)

// BuildImage builds a falco-driver-builder docker image for the given Falco Version and returns the built image's FQN.
func BuildImage(
	dockerClient *docker.Client,
	falcoVersion string,
) (string, error) {
	imageFQN := fmt.Sprintf("%s:%s", FalcoDriverBuilderRepository, falcoVersion)
	err := dockerClient.Build(&docker.BuildOpts{
		Dockerfile: FalcoDriverBuilderDockerfile,
		BuildArgs: map[string]*string{
			"FALCO_VERSION": docker.StrPtr(falcoVersion),
		},
		Tags: []string{imageFQN},
	})
	if err != nil {
		return "", fmt.Errorf("could not build %s: %w", imageFQN, err)
	}

	return imageFQN, nil
}

// GetProbeNameFromBuildOutput returns the built Falco probe name from the build output or an error if it could not be found.
func GetProbeNameFromBuildOutput(buildOutput string) (string, error) {
	reStr := strings.ReplaceAll(regexp.QuoteMeta(BuiltFalcoProbesDir)+`falco\_.*`, `/`, `\/`)
	re := regexp.MustCompile(reStr)
	probeMatch := re.FindString(buildOutput)
	if len(probeMatch) < 1 {
		return "", ErrCouldNotFindProbeNameInOutput
	}

	return filepath.Base(probeMatch), nil
}

// ExtractProbeFromVolume extracts the built Falco eBPF probe by its name from the given probe volume.
func ExtractProbeFromVolume(
	dockerClient *docker.Client,
	builtProbeVolume operatingsystem.Volume,
	probeName string,
) (io.Reader, error) {
	builtProbePath := filepath.Join(BuiltFalcoProbesDir, probeName)
	fileBytes, err := dockerClient.GetFileFromVolume(
		builtProbeVolume,
		BuiltFalcoProbesDir,
		builtProbePath,
	)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

// GetDriverVersion returns the Falco Driver Version for the given falco-driver-builder image.
func GetDriverVersion(dockerClient *docker.Client, image string) (string, error) {
	out, err := dockerClient.Run(&docker.RunOpts{
		Image:      image,
		Entrypoint: []string{"/bin/bash"},
		Cmd:        []string{"-c", "cat /usr/bin/falco-driver-loader | grep DRIVER_VERSION= | cut -f2 -d\\\""},
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// WriteProbeToFile writes the given probe bytes to the distribution structure determined by the falco driver version and probe name.
func WriteProbeToFile(falcoDriverVersion string, probeName string, probeReader io.Reader) (string, error) {
	outProbePath := filepath.Join("dist", falcoDriverVersion, probeName)
	if err := os.MkdirAll(filepath.Dir(outProbePath), os.ModePerm); err != nil {
		return "", err
	}

	f, err := os.Create(outProbePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := f.Chmod(0644); err != nil {
		return "", err
	}

	if _, err := io.Copy(f, probeReader); err != nil {
		return "", err
	}

	if err := f.Sync(); err != nil {
		return "", err
	}

	return outProbePath, nil
}
