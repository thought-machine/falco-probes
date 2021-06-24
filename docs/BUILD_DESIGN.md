# Falco eBPF Probe Building Design

[Falco](https://falco.org) is an open-source cloud-native runtime security tool that parses Linux system calls from the kernel and alerts when they match a user-defined set of rules.

## Overview

[Falco's architecture](https://falco.org/docs/getting-started/#falco-architecture) is [composed of](https://falco.org/docs/#what-are-the-components-of-falco):

- A _Falco Agent_: userspace daemon which processes these syscall events by:
  1. Matching on user-defined rules.
  2. Forwarding matches to user-defined outputs.
- A _Falco Driver_: kernel module/eBPF probe which collects Linux kernel syscall events.
- _Falco Configuration_ which includes rules to match on and where to forward rule-match events.

This design document is only concerned about the building of the _Falco Driver_ which is unique per Linux kernel.

## Goals

- Interoperable: We must be able to build Falco drivers for a wide-range of operating systems and their kernel versions (e.g. Amazon Linux 2, Google Container OS, Ubuntu, etc.).
- Scalable: Adding support for more operating systems or kernels must not increase the maintenance and build complexity exponentially. This project must not expect regular human maintenance.

## Background

### Why eBPF?

The _Falco Agent_ supports reading Linux syscall events from [3 types of *Falco Driver*s](https://falco.org/docs/event-sources/drivers/).

- **Kernel Module**
  - Advantages:
    - Works with the majority of kernel versions that support [DKMS](https://github.com/dell/dkms).
    - Can see root/kernelspace events.
  - Disadvantages:
    - [A bug in the kernel module can cause a system-wide outage](https://tldp.org/LDP/lkmpg/2.6/html/lkmpg.html#AEN1352).
    - [Unsupported by Google Container-Optimized OS](https://falco.org/docs/getting-started/third-party/production/#gke) and other operating systems with limited access to kernel features.
    - Requires the `--privileged` container flag.
- **eBPF Probe**
  - Advantages:
    - Supports [Google Container OS](https://falco.org/docs/getting-started/third-party/production/#gke) and other operating systems with limited access to kernel features.
    - Greater security control by restricting kernel logic, see [A thorough introduction to eBPF](https://lwn.net/Articles/740157/).
    - Can see root/kernelspace events.
  - Disadvantages:
    - Only supports Linux kernels >= 4.14.
    - Requires the `--privileged` container flag. In Linux kernels >= 5.8 `CAP_BPF` and `CAP_PERFMON` can be used instead.
- **Userspace instrumentation**
  - Advantages:
    - Does not require access to the kernel thus can run entirely unprivileged.
  - Disadvantages:
    - Currently (as of 06/2021) no officially supported implementation.
    - Cannot see root/kernelspace events.

The listing above favours the **eBPF Probe** driver as all of the limitations are primarily around support for legacy kernel versions. The [4.14 Kernel was released in 11/2017](https://lkml.org/lkml/2017/11/12/123) where the majority of actively supported cloud vendor provided operating systems that can be used with [Kubernetes](https://kubernetes.io) are using a Linux Kernel >= 4.14. With this in mind, eBPF's disadvantages can be considered moot for modern Kubernetes clusters (with a path for the `--privileged` container flag being mitigated in the future too).

### Building a Falco eBPF Probe

Falco provide the following methods to build _Falco Drivers_: 
* A bash script [`falco-driver-loader`](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader) which is used to build the Falco kernel module/eBPF probe at runtime. 
* [`driverkit`](https://github.com/falcosecurity/driverkit) which is a command-line tool in active development which FalcoSecurity use to produce their repository of Falco kernel probes.

`driverkit` is the de-facto way to build probes before runtime, however upon further investigation we found that it does not quite meet our needs:
- It currently does not support Google Container-Optimized OS (`cos`) which is non-trivial to add because there is no package repository and Google only publish kernel sources per build of `cos`. 
- This project is also used to pre-build Linux Kernel Modules in their https://download.falco.org/driver repository which is pushed via this [job](https://prow.falco.org/?job=build-drivers-amazonlinux2-periodic) on their CI system, [prow](https://github.com/falcosecurity/test-infra).
  - This job currently does not build eBPF probes, which we desire.
  - This job currently rebuilds every probe and re-uploads them which results in hash changes, which does not suit our static hash verification when fetching external assets.
  - Currently, supporting newer kernel versions requires a pull-request to the repository on GitHub, e.g. https://github.com/falcosecurity/test-infra/pull/419 which makes us dependent on FalcoSecurity's review processes.

With this in mind, we have favoured the `falco-driver-loader` method to give us broader Operating System support and attempt to resolve some of the current shortcomings of Falco's probe building infrastructure.

Reverse engineering the `falco-driver-loader` bash script yields the following inputs for building an eBPF probe:

- [`KERNEL_RELEASE`](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L508) which is the output of `uname -r`.
- [`KERNEL_VERSION`](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L514) which is the output of `uname -v` passed to a `sed` command which extracts the number after the `#`.
- [`ARCH`](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L506) which is the output of `uname -m`.
- [Target ID](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L103) which is obtained from `/etc/os-release`.
- [Kernel configuration](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L68) which is obtained from `/lib/modules/`.
- [Kernel sources](https://github.com/falcosecurity/falco/blob/0.28.1/scripts/falco-driver-loader#L387) which is obtained from `/usr/src/`.

Thus to build a Falco eBPF probe in Docker, we can:

1. Build a modified `falco-driver-loader` image (called `falco-driver-builder`) that allows us to mock these values by patching the `falco-driver-loader` script.
2. Obtain Kernel sources and configuration for your chosen kernel.
3. Mock the resolved _Target ID_ by mocking the `/etc/os-release` file.
4. Mock the _Kernel Release_ (output of `uname -r`) value.
5. Mock the _Kernel Version_ (output of `uname -v`) value.
6. Mock the _Kernel Machine_ (output of `uname -m`) value.
7. Build Probe using patched _falco-driver-loader_ script in _falco-driver-builder_ with mocked values, _Kernel sources_, _Kernel configuration_ and mocked _Target ID_.

This process is proved by an accompanying script which builds a Falco eBPF probe for `Amazon Linux 2` with the `4.14.232-176.381.amzn2` kernel which can be executed via:

```bash
# A list of Kernel packages for Amazon Linux 2 can be obtained by running:
# $ docker run --rm amazonlinux:2 yum --showduplicates list kernel-devel | tail -n+3 | awk '{ print $2 }'
$ bash ./docs/BUILD_DESIGN_assets/build-amazonlinux2-ebpf-probe.sh "4.14.232-176.381.amzn2"
```

## Design

In [Building a Falco eBPF Probe](#building-a-falco-ebpf-probe), we identified 6 inputs to the [`falco-driver-loader`](https://falco.org/docs/getting-started/installation/#install-driver) script. These are the requirements for building Falco eBPF probes for any Linux kernel. In order to support additional Operating Systems and their kernels in this project, we can define these 6 required inputs as functions within an _Interface_ to provide a layer of abstraction between different kernels.

```golang
// KernelPackage abstracts the implementation of resolving the required inputs
// for building a Falco eBPF probe per Kernel Package.
// The outputs of are not guaranteed to be unique, see "Operating Systems without Package Managers"
// below for an explanation.
// Note: This interface is provided as an example.
// It will be different in the implementation of this design
// to include scope such as error handling.
type KernelPackage interface {
    // GetKernelRelease returns the value to mock as the output of `uname -r`.
    GetKernelRelease() string
    // GetKernelVersion returns the value to mock as the output of `uname -v`.
    GetKernelVersion() string
    // GetKernelMachine returns the value to mock as the output of `uname -m`.
    GetKernelMachine() string
    // GetOSRelease returns the file contents to use as the mock of `/etc/os-release`.
    GetOSRelease() FileContents
    // GetKernelConfiguration returns the volume to mount as `/host/lib/modules/`.
    GetKernelConfiguration() Volume
    // GetKernelSources returns the volume to mount as `/usr/src/`.
    GetKernelSources() Volume
}

// e.g. for the `4.14.232-176.381.amzn2 kernel` on `Amazon Linux 2`:
// GetKernelRelease() returns "4.14.232-176.381.amzn2.x86_64".
// GetKernelVersion() returns "#1 SMP Wed May 19 00:31:54 UTC 2021".
// GetKernelMachine() returns "x86_64".
// GetOSRelease() returns the contents /etc/os-release file from the Amazon Linux 2 docker image.
// GetKernelConfiguration() returns the volume with `/lib/modules/` for the kernel after running `yum install -y kernel-...`.
// GetKernelSources() returns the volume with `/usr/src/` for the kernel after running `yum install -y kernel-...`.
```

The above interface does not cover step 2 of the [Building a Falco eBPF Probe](#building-a-falco-ebpf-probe) process i.e. _2. Obtain Kernel sources and configuration for your chosen kernel_. For the `4.14.232-176.381.amzn2` kernel on `Amazon Linux 2`, we obtained these by running the `yum -y install "kernel-devel-$KERNEL_PACKAGE" "kernel-$KERNEL_PACKAGE"` command. This command utilises the `yum` package manager which is specific to RHEL and its children of which Amazon Linux 2 is one of. However, this command doesn't work in other Operating Systems such as `Ubuntu Linux` or `Google Container OS` and thus needs abstraction.

Additionally, this project aims to build Falco eBPF Probes for Kernel Packages from different Operating Systems (i.e. _Interoperability_). To meet this goal without increasing maintenance complexity, we can programmatically retrieve a list of Kernel Packages for an Operating System. In the Amazon Linux 2 example, we achieved this by running the `yum --showduplicates list kernel-devel` command. Again, this command does not work in other Operating Systems such as `Ubuntu Linux` or `Google Container OS` and thus needs abstracting as well.

The _Interface_ below abstracts these 2 Operating System requirements as 2 functions.

```golang
// OperatingSystem abstracts the implementation of determining which
// kernel packages are available and the retrieval of them.
// Note: This interface is provided as an example.
// It will be different in the implementation of this design
// to include scope such as error handling.
type OperatingSystem interface {
    // GetKernelPackageNames returns a list of all available Kernel Package names.
    GetKernelPackageNames() []string
    // GetKernalPackageByName returns a hydrated KernelPackage for the given Kernel Package name.
    GetKernelPackageByName(name string) KernelPackage
}

// e.g. for the `4.14.232-176.381.amzn2 kernel` on `Amazon Linux 2`:
// GetKernelPackageNames() returns []string{"4.14.232-176.381.amzn2", ...}.
// GetKernelPackageByName("4.14.232-176.381.amzn2") returns the example KernelPackage above.
```

Note: _hydrated_ means that the values are retrieved, i.e. the `GetKernelPackageByName` function performs the fetching of _Kernel Sources_, _Kernel Configuration_, etc.

### Operating Systems without Package Managers

Not all Operating Systems feature a Package Manager (e.g. `yum`, `apt`, `pacman`, etc.) thus `KernelPackage` may be seen as misnamed in the context of those Operating Systems. An example of this is `Google Container OS` which features security measures such as immutable filesystems. In order to fit these types of Operating Systems, we can use their [`BuildID`](https://cloud.google.com/container-optimized-os/docs/release-notes) as a `KernelPackage` where multiple `KernelPackage`s may output the same Falco eBPF Probe in the likely event that `BuildID`s share Kernels.

### Building Falco eBPF Probes at Scale

Now that we have our _Interfaces_ which abstract the Operating Systems and their Kernels (_Interoperability_), we need to design our implementation for using these to build Falco eBPF probes at _Scale_.

The below binaries are currently separated to enforce separation of concerns, but it is plausible that these may be merged into a single binary in a future maturity.

#### `//cmd/build-falco-ebpf-probe`

In [Building a Falco eBPF Probe](#building-a-falco-ebpf-probe), we used Bash to orchestrate the building of an eBPF probe via [Docker](https://docker.com). Bash can be replaced by a [Go](https://golang.org/) binary which utilises the above _Interfaces_ to perform the build steps via the [Docker SDK](https://pkg.go.dev/github.com/docker/docker/client).

```bash
$ plz run //cmd/build-falco-ebpf-probe -- <operating-system> <kernel-package-name>
# $ plz run //cmd/build-falco-ebpf-probe -- amazonlinux2 4.14.232-176.381.amzn2
# Built eBPF probe for 4.14.232-176.381.amzn2 on amazonlinux2.
```

#### `//cmd/list-kernel-packages`

We will also require a Go binary which can list the available Kernel Packages for a given Operating System.

```bash
$ plz run //cmd/list-kernel-packages -- <operating-system>
# $ plz run //cmd/list-kernel-packages -- amazonlinux2
# 4.14.232-176.381.amzn2
# ...
```

This binary's output will be used in CI/CD to build Falco eBPF probes for all available kernels.

#### `//cmd/is-falco-ebpf-probe-uploaded`

This Go binary will be used in CI/CD to determine whether or not a Falco eBPF probe has already been uploaded to our _Probe Repository_.
The implementation of this depends on `docs/REPOSITORY_DESIGN.md`.

```bash
$ plz run //cmd/is-falco-ebpf-probe-uploaded -- <operating-system> <kernel-package-name>
# $ plz run //cmd/is-falco-ebpf-probe-uploaded -- amazonlinux2 4.14.232-176.381.amzn2
# (exit 0 - probe exists in repository)
# (exit 1 - probe does not exist in repository)
```

#### `Volume`s

In the `KernelPackage` _Interface_, we have referenced a `Volume` datatype. This is used to abstract from the implementation of different file storage mechanisms. For the first implementation of this design, we recommend that [Docker Volumes](https://docs.docker.com/storage/volumes/) are used.

#### GitHub Actions

GitHub Actions offers us a "free" and transparent way for us to build our eBPF probes as well as integrate with the GitHub Repository directly. There is an undocumented parallel worker limit of 256 on GitHub Actions, which means that we cannot build all available Kernels in parallel which also results in a very noisy GitHub Actions UI as there would be a "Job" for each Kernel.

Instead, we suggest that we use a GitHub Actions worker per Operating System, which falls well within our scaling needs as jobs can run for up to 2 hours.

In order to automatically build new eBPF probes, we propose to initially run these tools on a nightly cron-job. We will only build Falco eBPF probes for Kernels which do not already exist in the _Probe Repository_, so this will be quiet after the initial builds that populate the _Probe Repository_.

## Future Considerations

This design aims to be representative of what a 1st iteration of maturity for this project could look like. There are certainly further improvements that can be made and should be considered in future maturities.

### Operating System Implementation Specific Optimisations

- **Fetching Kernel Sources and Configuration**: Currently, we've demonstrated the use of Operating System Package Managers to fetch Kernel Sources and Configuration which is inefficient as the package managers fetch additional packages. We could improve this by directly fetching from the repositories via HTTP.
- **Listing Kernel Packages**: Currently, we've demonstrated the use of Operating System Package Managers to list available Kernel Packages which is inefficient as it requires running a Docker command to do so. We could improve this by directly fetching from the repositories via HTTP.

### Building Entirely in Please without Docker

Currently, we're advising to build entirely in Docker via the `docker.io/falcosecurity/falco-driver-loader` base image which comes with all the build dependencies required. However, this requires running and depending on Docker which can be inefficient, where building entirely with Please would be much cleaner.

### Automated verification of built eBPF Probes

It's possible that our build process may produce incompatible probes, we could build an E2E style test which tests our built eBPF probes against their respective kernels.


### Supply-chain security of built eBPF Probes

As we're building eBPF programs which run inside the Linux Kernel, it is desirable for us to provide a way for consumers of this project to verify that the probe they have downloaded was built by us.
