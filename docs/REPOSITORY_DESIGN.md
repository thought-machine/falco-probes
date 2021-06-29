# Falco eBPF Probe Repository Design

## Overview

In [BUILD_DESIGN.md](./BUILD_DESIGN.md), we outlined an approach to building Falco eBPF probes for a wide-range of operating systems and their kernel versions. This project aims to host these pre-compiled probes publicly, in a convenient and comprehensive repository; individual probes should be easily available to download, addressable by their runtime characteristics (driver version/OS/kernel version) and an extensive catalog of probes should be available for a range of Falco versions.

## Goals

We want to maintain a probe repository that satisfies the following goals:
- For public consumption: Probes should be readily and freely available to download in a manner that is simple and intuitive
- Low Overhead: The repository should be maintainable with little/no infrastructure overhead or cost to the project maintainers
- falco-driver-loader Compatible: Probes should be available at a URL that matches the pattern expected by falco-driver-loader
- Trustworthy: The automation from compilation to upload should be observable and verifiable to consumers, while the repository itself should be hosted on trusted, secure-by-default infrastructure

Availability and latency/performance of the repository are not prioritised, primarily because we should advise against consumers of this project from having our registry as a runtime dependency.

## Design

To achieve the aforementioned goals, we can leverage Github Releases as our probe repository, serving compiled probes to consumers as downloadable release assets, enabling us to keep the entirety of the source-code, compilation process and hosting of probes consolidated together on a single, trusted and widely used platform. 

Releases can be organised as follows:
- Each new version of the Falco driver can be associated to a single github release, named to match the semver of the driver version.
- Each compiled probe can be uploaded as an asset to the release of the driver version the probe was compiled against, named based on the runtime characteristics of the OS/kernel.

E.g
```
- 0.18.0
    - falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.225-169.362.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.181-142.260.amzn2.x86_64_1.o
- 0.17.1
    - falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.225-169.362.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.181-142.260.amzn2.x86_64_1.o
- 0.17.0
    - falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.225-169.362.amzn2.x86_64_1.o
    - falco_amazonlinux2_4.14.181-142.260.amzn2.x86_64_1.o
...
```

To download a particular probe based on the driver version/OS/kernel version, you would use the following URL
`https://github.com/thought-machine/falco-probes/releases/download/$DRIVER_VERSION/$PROBE_FILENAME.o`.

For example, to download the first probe listed above via curl you would simply need to:
`curl -L https://github.com/thought-machine/falco-probes/releases/download/0.18.0/falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o > falco_amazonlinux2_4.14.232-177.418.amzn2.x86_64_1.o`

As new OS/kernel version combinations become available, all prior releases can be updated in parallel to include assets for each newly compiled probe. 

There is no authentication required to download assets from Github, no rate-limiting or throttling applied to downloads, all at no cost to the maintainers.

### Attaching Assets to a Release
Probes can be added as assets to a release using the [gh-release](https://github.com/marketplace/actions/gh-release) action. This can be done from the same job that compiles the probe itself; we could simply mount the `/root/.falco` directory of the container to the current working directory of the job, and then run a 'release' step after the 'build' step, defined something like:
```
  - name: Release
    uses: softprops/action-gh-release@v1
    env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    with:
        tag_name: ${{DRIVER_VERSION}}
        files: |
            ${{PROBE_FILENAME}}.o # available in our pwd thanks to the volume
```

The first time this would run for a given DRIVER_VERSION, the release will be created; subsequent runs will simply edit the existing release to append more assets to the asset list.

The `secrets.GITHUB_TOKEN` variable is [automatically available from the action](https://docs.github.com/en/actions/reference/authentication-in-a-workflow#about-the-github_token-secret), so no additional configuration is required. There are no *documented* limits to the number of files you can associate with a release, so this process can be continually repeated for as long as a $DRIVER_VERSION is supported.

### Releases, Tags and Refs

Github Releases are associated to tags, tying the released assets to the source code at the time said assets were built; this is both the intent design and the user expectation with Releases. If we choose to continually update prior releases to add new assets, we will break this expectation; there will be no guarantee a particular probe was compiled against the source taken at the release/tag ref. Functionally, this won't impact our repository, but it will make it more difficult for us to triage compiled probes, as well as making it more difficult to verify, or reproduce, the compilation process for a given probe in a release. If possible, it would be nice to make individual probe compilation runs discoverable in the 'workflow filter' in the Actions page, simplifying the process of cross-referencing a built probe to the commit the compilation was performed against.

### Digests

With any stored asset, a digest should be produced to enable consumers to verify the asset has not been modified, maliciously or otherwise. Rather than accompanying each probe with a .sha256 digest file, which would double the number of files present in each release, we can simply print a digest from within the context of the compilation action, allowing us and any other consumers to verify the asset in the release matches the compiled version by reviewing the logs in the action.

## Limits or Restrictions on Releases

The only documented constraint associated with releases is a [max file-size constraint of 2GB](https://docs.github.com/en/github/administering-a-repository/releasing-projects-on-github/about-releases#storage-and-bandwidth-quotas). Since we will be producing large releases, and editing individual releases multiple times, we've [deployed a crack team of hamsters to stress-test releases](https://github.com/sHesl/release-loop/actions) by continually updating an existing release with new assets. If there are any undocumented restrictions around the number of edits you can apply to a release, or the number of files associated to a single release, they will find it!

## Future Considerations

### Discoverability
Releases, especially those with lots of assets, aren't easy to explore. Users who wish to check for the existence of a particular probe have several, unsatisfactory options:
- Use the [Releases API](https://docs.github.com/en/rest/reference/repos#list-releases) to iterate through past releases, and then again through the assets for a given release.
- Manually scrolling through the release page in the Github UI.
- Trying to curl down the probe, using the status code to infer it's existence.
A more user-friendly way of documenting/highlighting which probes have already been compiled could be beneficial. One approach might be to maintain a CHANGELOG type file in the repo that is automatically updated upon successful upload that contains information about the build (link to the action, copy of the digest etc).