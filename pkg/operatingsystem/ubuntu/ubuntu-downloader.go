package ubuntu

import (
	"fmt"
	"strings"

	"github.com/thought-machine/falco-probes/pkg/docker"
)

// UbuntuDownloaderDockerfile represents the contens of a Dockerfile to build the ubuntudownloader image which
// downloads DEB packages from Ubuntu's repositories. We add each supported Codename's sources to APT.
const UbuntuDownloaderDockerfile = `FROM ubuntu:latest
RUN export codenames="" \
	&& . /etc/os-release \
	&& for c in $codenames; do \
		cp /etc/apt/sources.list "/etc/apt/sources.list.d/${c}.list" \
		&& sed -i "s/${UBUNTU_CODENAME}/${c}/g" /etc/apt/sources.list.d/${c}.list; \
	done
`

// UbuntuDownloaderRepository is the repository to build the UbuntuDownloader image under.
const UbuntuDownloaderRepository = "docker.io/thoughtmachine/falco-ubuntudownloader"

// SupportedCodenames represents the supported versions of Ubuntu.
// These should be updated for new releases and removed when they reach End of Standard Support (https://wiki.ubuntu.com/Releases).
var SupportedCodenames = []string{
	"bionic", // 18.04 LTS, EoSS: April 2023
	"focal",  // 20.04 LTS, EoSS: April 2025
}

// BuildUbuntuDownloader builds the UbuntuDownloader docker image.
func BuildUbuntuDownloader(dockerClient *docker.Client) (string, error) {
	imageFQN := fmt.Sprintf("%s:latest", UbuntuDownloaderRepository)
	err := dockerClient.Build(&docker.BuildOpts{
		Dockerfile: strings.Replace(UbuntuDownloaderDockerfile, `codenames=""`, `codenames="`+strings.Join(SupportedCodenames, " ")+`"`, 1),
		Tags:       []string{imageFQN},
	})
	if err != nil {
		return "", fmt.Errorf("could not build %s: %w", imageFQN, err)
	}

	return imageFQN, nil
}
