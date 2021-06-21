#!/bin/bash
# This script demonstrates the building of an Amazon Linux 2 Falco eBPF probe.
set -Eeuox pipefail

# FALCO_VERSION is the version of Falco to compile the eBPF probe for.
FALCO_VERSION="0.28.1"
# The chosen kernel is defined as `KERNEL_PACKAGE`.
if [ -v 1 ]; then
    echo "A kernel package version is expected as an argument."
    exit 2
fi
KERNEL_PACKAGE="$1"
# FALCO_DRIVER_BUILD_IMAGE is the docker image tag for the patched version of falco-driver-loader
FALCO_DRIVER_BUILDER_IMAGE="falco-driver-loader:${FALCO_VERSION}"

# 1. Build a modified `falco-driver-loader` image (called `falco-driver-builder`) that supports inputs:
# - `UNAME_R` (mock output of `uname -r`)
# - `UNAME_V` (mock output of `uname -v`)
# - `UNAME_M` (mock output of `uname -m`)
# and sets the `KERNELDIR` to `/lib/modules/$KERNEL_RELEASE/build` which is used in Makefiles.
docker build \
    --tag "${FALCO_DRIVER_BUILDER_IMAGE}" \
    --build-arg "FALCO_VERSION=${FALCO_VERSION}" - \
    < docs/BUILD_DESIGN_assets/falco-driver-builder.Dockerfile

# 2. Obtain Kernel sources and configuration for the chosen kernel.
# Note: We use Docker Volumes to store the unpacked data.
usr_src_volume=$(docker volume create)
lib_modules_volume=$(docker volume create)
docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --volume "${lib_modules_volume}":/lib/modules/ \
    amazonlinux:2 yum -y install "kernel-devel-$KERNEL_PACKAGE" "kernel-$KERNEL_PACKAGE"

# 3. Mock the `/etc/os-release` file.
etc_volume=$(docker volume create)
docker run --rm \
    --volume "${etc_volume}":/host/etc/ \
    amazonlinux:2 cp /etc/os-release /host/etc/

# Find the *Kernel source path* to determine the *Kernel Release* and *Kernel Version* from.
KERNEL_SRC_PATH=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    find /usr/src/ -name "*$KERNEL_PACKAGE*" -type d
)

# 4. Mock the *Kernel Release* (`uname -r`) by running `make kernelrelease` in the $KERNEL_SRC_PATH.
UNAME_R=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "make kernelrelease" | tail -n1
)

# 5. Mock the *Kernel Version* (`uname -v`) by looking in the `./include/generated/compile.h` file in the $KERNEL_SRC_PATH.
UNAME_V=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_VERSION.* | cut -f2 -d\\\""
)

# 6. Mock the *Kernel Machine* (output of `uname -m`) by looking in the `./include/generated/compile.h` file in the $KERNEL_SRC_PATH.
UNAME_M=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    -w "$KERNEL_SRC_PATH" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "find /usr/src -name compile.h | grep 'generated/compile.h' | xargs grep -ho UTS_MACHINE.* | cut -f2 -d\\\""
)

# 7. Build Probe using patched *falco-driver-loader* script in *falco-driver-builder* with mocked values, *Kernel sources*, *Kernel configuration* and mocked *Target ID*.
docker run --rm \
    --env UNAME_V="$UNAME_V" \
    --env UNAME_R="$UNAME_R" \
    --env UNAME_M="$UNAME_M" \
    --volume "${usr_src_volume}":/host/usr/src/ \
    --volume "${lib_modules_volume}":/host/lib/modules/ \
    --volume "${etc_volume}":/host/etc/ \
    "${FALCO_DRIVER_BUILDER_IMAGE}"

# Clean up Docker volumes
docker volume rm "${usr_src_volume}" "${lib_modules_volume}" "${etc_volume}"
