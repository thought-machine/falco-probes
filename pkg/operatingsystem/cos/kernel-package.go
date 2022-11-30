package cos

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

const (
	// BusyBoxImage is the image to use for BusyBox, this is used as the smallest Docker image that we can create a container from.
	// We cannot use scratch unless we build our own blank image from scratch, which adds more complexity.
	BusyBoxImage = "docker.io/library/busybox:1.33.1"

	kernelReleasePattern       = `^([0-9]+\.){2}[0-9]+$`
	urlCosKernelConfigTemplate = "https://cos.googlesource.com/third_party/kernel/+/%s/arch/x86/configs/%s_defconfig?format=TEXT"
	urlCosToolsTemplate        = "https://storage.googleapis.com/cos-tools/%s/%s"
)

var (
	arches = [...]string{"lakitu", "x86_64"}
	log    = logging.Logger
)

// NewKernelPackage returns a new hydrated example implementation operatingsystem.KernelPackage.
func NewKernelPackage(dockerClient *docker.Client, name string) (*operatingsystem.KernelPackage, error) {
	kP := &operatingsystem.KernelPackage{
		OperatingSystem: "cos",
		Name:            name,
	}

	version, err := ParseVersion(name)
	if err != nil {
		return nil, err
	}

	if err := addKernelReleaseAndVersionAndMachine(dockerClient, kP, version); err != nil {
		return nil, err
	}

	if err := addSourcesAndConfiguration(dockerClient, kP, version); err != nil {
		return nil, err
	}

	if err := addOSRelease(dockerClient, kP, version); err != nil {
		return nil, err
	}

	return kP, nil
}

// Falco doesn't require the Google COS sources so don't bother to fetch them and just create an empty volume for the
// interface.
func addSourcesAndConfiguration(dockerClient *docker.Client, kp *operatingsystem.KernelPackage, version *Version) error {
	kernelCommit, err := readKernelCommit(version.BuildID)
	if err != nil {
		return err
	}

	encodedKernelConfig, err := readKernelConfig(version.BuildID, kernelCommit)
	if err != nil {
		return err
	}

	kernelConfig, err := decodeKernelConfig(version.BuildID, kernelCommit, encodedKernelConfig)
	if err != nil {
		return err
	}

	commandTemplate := "mkdir -p /usr/src/kernels && mkdir -p '/lib/modules/%s+'"

	kp.KernelConfiguration = dockerClient.MustCreateVolume()
	kp.KernelSources = dockerClient.MustCreateVolume()
	_, err = dockerClient.Run(
		&docker.RunOpts{
			Image:      BusyBoxImage,
			Entrypoint: []string{"/bin/sh"},
			Cmd:        []string{"-c", fmt.Sprintf(commandTemplate, kp.KernelRelease)},
			Volumes: map[operatingsystem.Volume]string{
				kp.KernelSources:       "/usr/src/",
				kp.KernelConfiguration: "/lib/modules/",
			},
		},
	)
	if err != nil {
		return err
	}

	err = dockerClient.WriteFileToVolume(kp.KernelConfiguration, "/lib/modules/", fmt.Sprintf("/lib/modules/%s+/config", kp.KernelRelease), kernelConfig)
	if err != nil {
		return err
	}

	return nil
}

func addOSRelease(dockerClient *docker.Client, kp *operatingsystem.KernelPackage, version *Version) error {
	osReleaseTemplate := `
ID=cos
NAME="Container-Optimized OS"
PRETTY_NAME="Container-Optimized OS from Google"
VERSION=%d
VERSION_ID=%d
BUILD_ID=%s
`
	osRelease := fmt.Sprintf(osReleaseTemplate, version.Milestone, version.Milestone, version.BuildID)

	osReleaseVol := dockerClient.MustCreateVolume()

	err := dockerClient.WriteFileToVolume(osReleaseVol, "/host/etc/", "/host/etc/os-release", osRelease)
	if err != nil {
		return err
	}

	fileReader, err := dockerClient.GetFileFromVolume(osReleaseVol, "/host/etc/", "/host/etc/os-release")
	if err != nil {
		return err
	}

	fileContents, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return err
	}

	kp.OSRelease = operatingsystem.FileContents(fileContents)

	return nil
}

func addKernelReleaseAndVersionAndMachine(dockerClient *docker.Client, kp *operatingsystem.KernelPackage, version *Version) error {
	kernelHeaders, err := readKernelHeaders(version.BuildID)
	if err != nil {
		return err
	}
	defer kernelHeaders.Close()
	return extractKernelDetails(version.BuildID, kernelHeaders, kp)
}

func readKernelHeaders(buildID string) (io.ReadCloser, error) {
	resp, err := http.Get(fmt.Sprintf(urlCosToolsTemplate, buildID, "kernel-headers.tgz"))
	if err != nil {
		return nil, fmt.Errorf("could not get kernel headers for build id %s: %w", buildID, err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("could not get 2XX response for kernel headers for build id %s: %s", buildID, resp.Status)
	}

	return resp.Body, nil
}

func extractKernelDetails(buildID string, kernelHeaders io.Reader, kp *operatingsystem.KernelPackage) error {
	decompressedKernelHeaders, err := gzip.NewReader(kernelHeaders)
	if err != nil {
		return fmt.Errorf("could not decompress kernel headers for build id %s: %w", buildID, err)
	}

	extractedKernelHeaders := tar.NewReader(decompressedKernelHeaders)

	for true {
		header, err := extractedKernelHeaders.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("could not get next file header in kernel headers archive for build id %s: %w", buildID, err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// Read kernel release from directory e.g. ./usr/src/linux-headers-5.15.73+/ -> 5.15.73
			if kp.KernelRelease != "" {
				continue
			}
			components := strings.SplitN(header.Name, "linux-headers-", 2)
			if len(components) < 2 {
				continue
			}
			components = strings.SplitN(components[1], "+", 2)
			if len(components) < 2 {
				continue
			}
			// Validate kernel release.
			re := regexp.MustCompile(kernelReleasePattern)
			if re.FindString(components[0]) == "" {
				return fmt.Errorf("could not validate kernel release '%s' against pattern '%s' in kernel headers archive for build id %s", components[0], kernelReleasePattern, buildID)
			}
			kp.KernelRelease = components[0]
		case tar.TypeReg:
			// Read kernel version and machine from generated/compile.h
			if kp.KernelVersion != "" && kp.KernelMachine != "" {
				continue
			}
			if !strings.Contains(header.Name, "generated/compile.h") {
				continue
			}
			scanner := bufio.NewScanner(extractedKernelHeaders)
			for scanner.Scan() {
				line := scanner.Text()
				switch {
				case strings.Contains(line, "UTS_MACHINE"):
					kp.KernelMachine = strings.SplitN(line, "\"", 3)[1]
				case strings.Contains(line, "UTS_VERSION"):
					kp.KernelVersion = strings.SplitN(line, "\"", 3)[1]
				default:
					continue
				}
			}
			err := scanner.Err()
			if err != nil {
				return fmt.Errorf("could not read lines of generate in kernel headers archive for build id %s: %s in %s", buildID, header.Typeflag, header.Name)
			}
		}

		if kp.KernelRelease != "" && kp.KernelVersion != "" && kp.KernelMachine != "" {
			break
		}
	}

	return nil
}

func readKernelCommit(buildID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(urlCosToolsTemplate, buildID, "kernel_commit"))
	if err != nil {
		return "", fmt.Errorf("could not get kernel commit for build id %s: %w", buildID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("could not get 2XX response for kernel commit for build id %s: %s", buildID, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read kernel commit for build id %s: %w", buildID, err)
	}

	return strings.TrimSuffix(string(body), "\n"), nil
}

func readKernelConfig(buildID string, kernelCommit string) (string, error) {
	body := make([]byte, 0)

	archLastIndex := len(arches) - 1
	for i, arch := range arches {
		resp, err := http.Get(fmt.Sprintf(urlCosKernelConfigTemplate, kernelCommit, arch))
		if err != nil {
			return "", fmt.Errorf("could not get kernel config for build id %s (kernel commit %s): %w", buildID, kernelCommit, err)
		}

		// If the config for the preferred architecture is not found, retry with the next one.
		if resp.StatusCode == 404 && i < archLastIndex {
			log.Warn().
				Str("build_id", buildID).
				Str("kernel_commit", kernelCommit).
				Str("architecture", arch).
				Msg("could not find config for")
			continue
		}

		if resp.StatusCode > 299 {
			return "", fmt.Errorf("could not get 2XX response for kernel config for build id %s (kernel commit %s): %s", buildID, kernelCommit, resp.Status)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("could not read kernel commit for build id %s: %w", buildID, err)
		}

		break
	}

	return string(body), nil
}

func decodeKernelConfig(buildID string, kernelCommit string, encodedKernelConfig string) (string, error) {
	decodedKernelConfig, err := base64.StdEncoding.DecodeString(encodedKernelConfig)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode kernel config for build id %s (kernel commit %s): %w", buildID, kernelCommit, err)
	}

	return string(decodedKernelConfig), nil
}
