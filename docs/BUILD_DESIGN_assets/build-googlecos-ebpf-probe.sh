#!/bin/bash
# This script demonstrates the building of an Google COS Falco eBPF probe.
set -Eeuox pipefail

# FALCO_VERSION is the version of Falco to compile the eBPF probe for.
FALCO_VERSION="0.33.0"
# UBUNTU_VERSION is the version of Ubuntu to use as the build environment.
UBUNTU_VERSION="22.04"
# The chosen kernel is defined as `IMAGE_NAME`.
if [ -v $1 ]; then
    echo "An image name is expected as an argument."
    exit 2
fi
# Examples:
# cos-101-17162-40-1 (BUILD_ID: 17162.40.1)
# cos-101-17162-40-20 (BUILD_ID: 17162.40.20)
IMAGE_NAME="$1"
# See https://cloud.google.com/container-optimized-os/docs/concepts/versioning
# VERSION is IMAGE_NAME with alpha components like `cos-dev-...` stripped from the front.
VERSION="$(sed -e "s/^\([a-z]\+-\)*//g" <<< "${IMAGE_NAME}")"
# MILESTONE is the first numeric component.
MILESTONE="${VERSION%%-*}"
# BUILD_NUMBER is the rest.
BUILD_NUMBER="${VERSION#*-}"
# BUILD_ID is BUILD_NUMBER with `.` instead of `-`.
BUILD_ID="${BUILD_NUMBER//-/.}"
# FALCO_DRIVER_BUILD_IMAGE is the docker image tag for the patched version of falco-driver-loader
FALCO_DRIVER_BUILDER_IMAGE="falco-driver-builder:${FALCO_VERSION}"

# 1. Build a modified `falco-driver-loader` image (called `falco-driver-builder`) that supports inputs:
# - `UNAME_R` (mock output of `uname -r`)
# - `UNAME_V` (mock output of `uname -v`)
# - `UNAME_M` (mock output of `uname -m`)
# and sets the `KERNELDIR` to `/lib/modules/$KERNEL_RELEASE/build` which is used in Makefiles.
docker build \
    --tag "${FALCO_DRIVER_BUILDER_IMAGE}" \
    --build-arg "FALCO_VERSION=${FALCO_VERSION}" \
    --build-arg "UBUNTU_VERSION=${UBUNTU_VERSION}" \
    - \
    < docs/BUILD_DESIGN_assets/falco-driver-builder.Dockerfile

# 2. Obtain Kernel headers for the chosen kernel so we can get the kernel version.
# Notes:
# - We use Docker Volumes to store the unpacked data.
# - falco-driver-loader downloads the kernel headers separately so this is not used for the probe build.
usr_src_volume=$(docker volume create)
docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "\
      rm -rf /usr/src/* && \
      curl -sSL 'https://storage.googleapis.com/cos-tools/${BUILD_ID}/kernel-headers.tgz' | \
        tar xz"

# 3. Mock the `/etc/os-release` file.
etc_volume=$(docker volume create)
docker run --rm \
    --volume "${etc_volume}":/host/etc/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "printf 'ID=cos\nBUILD_ID=%s\nVERSION_ID=%s\n' '${BUILD_ID}' '${MILESTONE}' > /host/etc/os-release"

# 4. Find the *Kernel source path* to determine the *Kernel Release* and *Kernel Version* from.
KERNEL_SRC_PATH=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "cd /usr/src/* && pwd"
)

# 5. Mock the *Kernel Release* (`uname -r`) by running `make kernelrelease` in the $KERNEL_SRC_PATH.
UNAME_R=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "make kernelrelease" | tail -n1
)

# 6. Mock the *Kernel Version* (`uname -v`) by looking in the `./include/generated/compile.h` file in the $KERNEL_SRC_PATH.
UNAME_V=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_VERSION.* | cut -f2 -d\\\""
)

# 7. Mock the *Kernel Machine* (output of `uname -m`) by looking in the `./include/generated/compile.h` file in the $KERNEL_SRC_PATH.
UNAME_M=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_MACHINE.* | cut -f2 -d\\\""
)

# 8. Obtain kernel configuration for the chosen kernel.
lib_modules_volume=$(docker volume create)
docker run --rm \
    --volume "${lib_modules_volume}":/lib/modules/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "\
      mkdir -p /lib/modules/${UNAME_R} && \
      curl -sSL https://storage.googleapis.com/cos-tools/${BUILD_ID}/kernel_commit | \
        xargs -I {} curl -sSL 'https://cos.googlesource.com/third_party/kernel/+/{}/arch/x86/configs/lakitu_defconfig?format=TEXT' | \
        base64 -d > /lib/modules/${UNAME_R}/config"

# 9. Build Probe using patched *falco-driver-loader* script in *falco-driver-builder* with mocked values, *Kernel sources*, *Kernel configuration* and mocked *Target ID*.
docker run --rm \
    --env UNAME_V="${UNAME_V}" \
    --env UNAME_R="${UNAME_R}" \
    --env UNAME_M="${UNAME_M}" \
    --env BUILD_ID="${BUILD_ID}" \
    --env HOST_ROOT="/host" \
    --volume "${usr_src_volume}":"/host/usr/src/" \
    --volume "${lib_modules_volume}":"/lib/modules/" \
    --volume "${etc_volume}":"/host/etc/" \
    "${FALCO_DRIVER_BUILDER_IMAGE}"

## Clean up Docker volumes
docker volume rm "${usr_src_volume}" "${lib_modules_volume}" "${etc_volume}"
