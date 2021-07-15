package amazonlinux2

import (
	"fmt"

	"github.com/thought-machine/falco-probes/pkg/docker"
)

// YumDownloaderDockerfile represents the contents of a Dockerfile to build the yumdownloader image which downloads RPM packages from AmazonLinux 2's repositories.
const YumDownloaderDockerfile = `FROM amazonlinux:2
RUN yum install -y yum-utils
`

// YumDownloaderRepository is the repository to build the yumdownloader image under.
const YumDownloaderRepository = "docker.io/thoughtmachine/falco-yumdownloader"

// BuildYumDownloader builds the yumdownloader docker image.
func BuildYumDownloader(dockerClient *docker.Client) (string, error) {
	imageFQN := fmt.Sprintf("%s:latest", YumDownloaderRepository)
	err := dockerClient.Build(&docker.BuildOpts{
		Dockerfile: YumDownloaderDockerfile,
		Tags:       []string{imageFQN},
	})
	if err != nil {
		return "", fmt.Errorf("could not build %s: %w", imageFQN, err)
	}

	return imageFQN, nil
}
