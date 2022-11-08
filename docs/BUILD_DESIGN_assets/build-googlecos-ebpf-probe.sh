#!/bin/bash
# This script demonstrates the building of an Google COS Falco eBPF probe.
#set -Eeuox pipefail
set -x

# FALCO_VERSION is the version of Falco to compile the eBPF probe for.
FALCO_VERSION="0.28.1"
# GENTOO_VERSION is the version of Gentoo to create the COS toolchain with.
GENTOO_VERSION="20221031"
# The chosen kernel is defined as `IMAGE_NAME`.
#if [ -v 1 ]; then
#    echo "An image name is expected as an argument."
#    exit 2
#fi
# Examples:
# cos-101-17162-40-1 (BUILD_ID: 17162.40.1)
# cos-101-17162-40-20 (BUILD_ID: 17162.40.20)
IMAGE_NAME="$1"
IMAGE_NAME="cos-101-17162-40-1"
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
FALCO_DRIVER_BUILDER_IMAGE="falco-driver-loader:${FALCO_VERSION}"
# GENTOO_IMAGE is the docker image tag for the gentoo
GENTOO_IMAGE="gentoo:${GENTOO_VERSION}"
# COS repo url path endings (https://cos.googlesource.com/third_party/overlays/%s) and names.
COS_REPO_URLS="chromiumos-overlay eclass-overlay portage-stable"
COS_REPO_NAMES="chromiumos eclass-overlay portage-stable"
# Chroot directory to run falco-driver-loader in
CHROOT_DIR="/toolchain"

# 1. Build a modified `falco-driver-loader` image (called `falco-driver-builder`) that supports inputs:
# - `UNAME_R` (mock output of `uname -r`)
# - `UNAME_V` (mock output of `uname -v`)
# - `UNAME_M` (mock output of `uname -m`)
# and sets the `KERNELDIR` to `/lib/modules/$KERNEL_RELEASE/build` which is used in Makefiles.
docker build \
    --tag "${FALCO_DRIVER_BUILDER_IMAGE}" \
    --build-arg "FALCO_VERSION=${FALCO_VERSION}" \
    - \
    < docs/BUILD_DESIGN_assets/falco-driver-builder.Dockerfile

# 2. Build a modified `gentoo` image which is used to populate the COS toolchain, bash, coreutils, and sed.
docker build \
    --tag "${GENTOO_IMAGE}" \
    --build-arg "GENTOO_VERSION=${GENTOO_VERSION}" \
    --build-arg "COS_REPO_URLS=${COS_REPO_URLS}" \
    --build-arg "COS_REPO_NAMES=${COS_REPO_NAMES}" \
    - \
    < docs/BUILD_DESIGN_assets/gentoo.Dockerfile

# 3. Obtain Kernel headers for the chosen kernel.
# Note: We use Docker Volumes to store the unpacked data.
usr_src_volume=$(docker volume create)
lib_modules_volume=$(docker volume create)
docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --volume "${lib_modules_volume}":/lib/modules/ \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "\
      rm -rf /usr/src/* && \
      wget -qO- 'https://storage.googleapis.com/cos-tools/${BUILD_ID}/kernel-headers.tgz' | \
        tar xvz"

# 3. Mock the `/etc/os-release` file.
etc_volume=$(docker volume create)
docker run --rm \
    --volume "${etc_volume}":/host/etc/ \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "printf 'ID=cos\nBUILD_ID=%s\nVERSION_ID=%s\n' '${BUILD_ID}' '${MILESTONE}' > /host/etc/os-release"

# Find the *Kernel source path* to determine the *Kernel Release* and *Kernel Version* from.
KERNEL_SRC_PATH=$(docker run --rm \
    --volume "${usr_src_volume}":/usr/src/ \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "cd /usr/src/* && pwd"
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

# 7. Obtain kernel configuration and srcs for the chosen kernel.
docker run --rm \
    --volume "${lib_modules_volume}":/lib/modules/ \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "\
      mkdir -p /lib/modules//${UNAME_R} && \
      wget -qO- https://storage.googleapis.com/cos-tools/${BUILD_ID}/kernel_commit | \
        xargs -I {} wget -qO- 'https://cos.googlesource.com/third_party/kernel/+/{}/arch/x86/configs/lakitu_defconfig?format=TEXT' | \
        base64 -d > /lib/modules/${UNAME_R}/config && \
      wget -qO- 'https://storage.googleapis.com/cos-tools/${BUILD_ID}/kernel-src.tar.gz' | \
        tar xvz -C /lib/modules/${UNAME_R}"

# 8. Obtain the package versions for the chosen kernel.
PACKAGES=$(docker run --rm \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "\
      wget -qO- https://storage.googleapis.com/cos-tools/${BUILD_ID}/cos-package-info.json | \
        jq -r '\
          .[][] | \
          select(.name | match(\"^(bash|coreutils|curl|findutils|gawk|grep|gzip|make|patch|sed|tar)\$\")) | \
          .category+\"/\"+.name+\"-\"+.ebuild_version' | \
        tr '\n' ' '"
)

# 9. Obtain the package repo refs for the chosen kernel.
COS_REPO_REFS=$(docker run --rm \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "\
      URLS=(\${COS_REPO_URLS}) && \
      wget -qO- https://storage.googleapis.com/cos-tools/${BUILD_ID}/manifest.xml > /tmp/manifest.xml && \
      for n in \"\${URLS[@]}\"; do \
        grep \"src/third_party/\$n\" /tmp/manifest.xml | grep -o 'path=\"[^\"]*\" revision=\"[^\"]*\"'; done | \
      sort | cut -d '\"' -f 4 | tr '\n' ' '"
)

# 9. Obtain the packages and toolchain for the chosen kernel and install into the chroot.
chroot_volume=$(docker volume create)
docker run --rm \
    --volume "${chroot_volume}":"${CHROOT_DIR}/" \
    "${GENTOO_IMAGE}" \
    /bin/bash -c "\
      NAMES=(${COS_REPO_NAMES}) && \
      REFS=(${COS_REPO_REFS})
      for ((i=0; i<\${#NAMES[@]}; i++)); do \
        git -C /var/db/repos/\${NAMES[\$i]} -c advice.detachedHead=false checkout \${REFS[\$i]}; done && \
      PACKAGES=(${PACKAGES}) && \
      for d in /dev /lib64; do mkdir -p \"${CHROOT_DIR}\${d}\"; done && \
      { test -e \"${CHROOT_DIR}/dev/null\" || mknod -m=666 \"${CHROOT_DIR}/dev/null\" c 1 3 ; } && \
      { test -e \"${CHROOT_DIR}/dev/random\" || mknod -m=666 \"${CHROOT_DIR}/dev/random\" c 1 8 ; } && \
      for f in /lib64/ld-linux-x86-64.so.2 /lib64/libc.so.6 /lib64/libm.so.6; do \
        cp -a --parents \${f} \"${CHROOT_DIR}\"; done && \
      for p in \"\${PACKAGES[@]}\"; do ROOT=\"${CHROOT_DIR}\" emerge -qv =\${p}; done && \
      curl -sSL \"https://storage.googleapis.com/cos-tools/${BUILD_ID}/toolchain.tar.xz\" | \
        tar xvJ -C \"${CHROOT_DIR}\" && \
      ln -fs /usr/bin/clang \"${CHROOT_DIR}/bin/gcc\""

# 10. Copy resolv.conf, falco-driver-loader and falco src into the chroot dir.
docker run --rm \
    --volume "${chroot_volume}":"${CHROOT_DIR}/" \
    --entrypoint="" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" \
    /bin/bash -c "\
      for f in /etc/resolv.conf /usr/bin/falco-driver-loader /usr/src/falco-*; do \
        cp -a --parents \${f} \"${CHROOT_DIR}\"; done"

# 11. Build Probe in the chroot using patched *falco-driver-loader* script in *falco-driver-builder* with mocked values, *Kernel sources*, *Kernel configuration* and mocked *Target ID*.
# apt-get update ; apt-get --fix-broken -y install; apt-get install -y less vim
docker run --rm \
    --env UNAME_V="${UNAME_V}" \
    --env UNAME_R="${UNAME_R}" \
    --env UNAME_M="${UNAME_M}" \
    --env BUILD_ID="${BUILD_ID}" \
    --env CHROOT_DIR="${CHROOT_DIR}" \
    --env HOST_ROOT="/host" \
    --env CC="/usr/bin/clang" \
    --env CLANG="/usr/bin/clang" \
    --env LLC="/bin/llc" \
    --volume "${chroot_volume}":"${CHROOT_DIR}" \
    --volume "${usr_src_volume}":"${CHROOT_DIR}/host/usr/src/" \
    --volume "${lib_modules_volume}":"${CHROOT_DIR}/lib/modules/" \
    --volume "${etc_volume}":"${CHROOT_DIR}/host/etc/" \
    "${FALCO_DRIVER_BUILDER_IMAGE}" 2>&1 | tee /tmp/falco-gcp.log

# Clean up Docker volumes
docker volume rm "${usr_src_volume}" "${lib_modules_volume}" "${etc_volume}" "${chroot_volume}"
