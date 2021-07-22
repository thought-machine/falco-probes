# [Falco](https://github.com/falcosecurity/falco) Probes

This project automates the building and mirroring of eBPF kernel probes for use by [Falco](https://github.com/falcosecurity/falco) as an [*eBPF probe*](https://falco.org/blog/choosing-a-driver/#ebpf-probe) driver to consume system call information which feeds its runtime threat detection abilities.

## Features
Falco Security provide [falcosecurity/driverkit](https://github.com/falcosecurity/driverkit) for building Falco drivers and [download.falco.org/driver](https://download.falco.org/driver) which mirrors only Falco [*kernel module*](https://falco.org/blog/choosing-a-driver/#kernel-module) drivers built from [falcosecurity/test-infra](https://github.com/falcosecurity/test-infra). This didn't quite meet all of our requirements, thus the birth of this project with the following features:

 * Falco eBPF probe drivers are built and available to [download](https://github.com/thought-machine/falco-probes/releases).
 * Mirrored Falco eBPF probe driver hashes do not change as we only upload probes once.
 * Probes for new kernel versions are automatically built as we dynamically obtain a list of available kernel versions at build time.

## Getting Started

This project mirrors built Falco eBPF probes to GitHub releases, where they are organised per Falco Driver Version (see [docs/REPOSITORY_DESIGN.md](./docs/REPOSITORY_DESIGN.md) for more information.). 

To obtain an eBPF kernel probe, you can:

1. Determine the Falco Driver version that your version of Falco is using:
```bash
FALCO_VERSION=0.29.1
docker run --rm --entrypoint="" \
    docker.io/falcosecurity/falco:$FALCO_VERSION \
    cat /usr/bin/falco-driver-loader \
    | grep DRIVER_VERSION= \
    | cut -f2 -d\"
# 17f5df52a7d9ed6bb12d3b1768460def8439936d
```
2. Go to the [Releases](https://github.com/thought-machine/falco-probes/releases) and find the name which matches your Falco Driver Version. You can then download the eBPF probes you want from there.

Below is a scripted example to download probes:
```bash
FALCO_VERSION=0.29.1
PROBE_NAME="falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o"

FALCO_DRIVER_VERSION=$(docker run --rm --entrypoint="" "docker.io/falcosecurity/falco:${FALCO_VERSION}" cat /usr/bin/falco-driver-loader | grep DRIVER_VERSION= | cut -f2 -d\")
# truncate driver version to 8 characters to get the release tag.
RELEASE_TAG=$(printf "%.8s\n" "${FALCO_DRIVER_VERSION}")

curl -LO "https://github.com/thought-machine/falco-probes/releases/download/${RELEASE_TAG}/${PROBE_NAME}"
``` 


## Supported Operating Systems

* Amazon Linux 2 (`amazonlinux2`)

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to add support for additional operating systems.
