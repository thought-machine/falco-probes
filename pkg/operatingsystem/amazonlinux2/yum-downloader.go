package amazonlinux2

import (
	"fmt"

	"github.com/thought-machine/falco-probes/pkg/docker"
)

// YumDownloaderDockerfile represents the contents of a Dockerfile to build the yumdownloader image which
// downloads RPM packages from AmazonLinux 2's repositories. We enable additional kernels available
// via amazon-linux-extras (when we know they will compile).
const YumDownloaderDockerfile = `FROM amazonlinux:2
RUN yum install -y yum-utils 
RUN export REPOS="kernel-5.4" \
	&& for r in $REPOS; do \
		amazon-linux-extras enable $r && \
		# disable to allow us to obtain repositories which overlap.
		amazon-linux-extras disable $r; \
	done \
	# enable all repositories.
	&& sed -i 's/enabled = 0/enabled = 1/g' /etc/yum.repos.d/amzn2-extras.repo
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
