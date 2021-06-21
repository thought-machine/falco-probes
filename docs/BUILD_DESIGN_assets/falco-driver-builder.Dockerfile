ARG FALCO_VERSION=0.28.1
FROM "docker.io/falcosecurity/falco-driver-loader:${FALCO_VERSION}"

ENV FALCO_DRIVER_LOADER_PATH="/usr/bin/falco-driver-loader"

SHELL ["/bin/bash", "-c"]

RUN set -Eeuxo pipefail; \
    # Disable downloading from Falco driver repository.
    sed -i 's/ENABLE_DOWNLOAD=.*/ENABLE_DOWNLOAD=""/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Enable compilation of eBPF driver.
    sed -i 's/ENABLE_COMPILE=.*/ENABLE_COMPILE="yes"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/DRIVER=.*/DRIVER="bpf"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Echo the KERNEL_RELEASE, KERNEL_VERSION and ARCH
    sed -i -e '/^KERNEL_RELEASE=.*/a\' -e 'echo "KERNEL_RELEASE: $KERNEL_RELEASE"' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i -e '/^KERNEL_VERSION=.*/a\' -e 'echo "KERNEL_VERSION: $KERNEL_VERSION"' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i -e '/^ARCH=.*/a\' -e 'echo "ARCH: $ARCH"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Set the KERNELDIR from the KERNEL_RELEASE. This is used in kernel Makefiles.
    sed -i -e '/^KERNEL_RELEASE=.*/a\' -e 'export KERNELDIR="/lib/modules/$KERNEL_RELEASE/build"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Allow setting the outputs of `uname` via UNAME_* environment variables.
    sed -i 's/uname -r/echo "${UNAME_R:-$(uname -r)}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -v/echo "${UNAME_V:-$(uname -v)}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -m/echo "${UNAME_M:-$(uname -m)}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    echo "Done!"
