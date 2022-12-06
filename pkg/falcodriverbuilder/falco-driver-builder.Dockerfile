ARG FALCO_VERSION
ARG UBUNTU_VERSION

FROM "docker.io/falcosecurity/falco-driver-loader:${FALCO_VERSION}" as falco-driver-loader

ENV FALCO_DRIVER_LOADER_PATH="/usr/bin/falco-driver-loader"

SHELL ["/bin/bash", "-c"]

RUN set -Eeuxo pipefail; \
    # Determine DRIVERS_REPO, DRIVER_VERSION, DRIVER_NAME from existing falco-driver-loader to persist them in the 0.33.0 falco-driver-loader script.
    export DRIVERS_REPO=$(grep ^DRIVERS_REPO "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    export DRIVER_VERSION=$(grep ^DRIVER_VERSION "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    export DRIVER_NAME=$(grep ^DRIVER_NAME "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    # Use falco-driver-loader from 0.33.0 (the patches below work with that script.)
    curl -L https://raw.githubusercontent.com/falcosecurity/falco/0.33.0/scripts/falco-driver-loader \
    -o /usr/bin/falco-driver-loader && chmod +x /usr/bin/falco-driver-loader && \
    # Add KBUILD_MODNAME if it doesn't already exist.
    { grep -q "KBUILD_MODNAME" "/usr/src/falco-${DRIVER_VERSION}/driver_config.h" || \
        printf '\n#ifndef KBUILD_MODNAME\n#define KBUILD_MODNAME "falco"\n#endif' >> \
            "/usr/src/falco-${DRIVER_VERSION}/driver_config.h"; \
    } && \
    # Add an always-y rule in Falco's bpf Makefile for newer kernels if it doesn't already exist.
    { grep -q "always-y" "/usr/src/falco-${DRIVER_VERSION}/bpf/Makefile" || \
        sed -i -e '/^always .*/a\' -e 'always-y += probe.o' "/usr/src/falco-${DRIVER_VERSION}/bpf/Makefile"; \
    } && \
    # Set DRIVERS_REPO to match existing falco-driver-loader script.
    sed -i '/^\s*DRIVERS_REPO=/d' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i "2 i DRIVERS_REPO=\"${DRIVERS_REPO}\"" "${FALCO_DRIVER_LOADER_PATH}" && \
    # Set DRIVER_VERSION to match existing falco-driver-loader script.
    sed -i '/^\s*DRIVER_VERSION=/d' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i "2 i DRIVER_VERSION=\"${DRIVER_VERSION}\"" "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i '3 i echo "DRIVER_VERSION: $DRIVER_VERSION"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Set DRIVER_NAME to match existing falco-driver-loader script.
    sed -i '/^\s*DRIVER_NAME=/d' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i "2 i DRIVER_NAME=\"${DRIVER_NAME}\"" "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i '3 i echo "DRIVER_NAME: $DRIVER_NAME"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Set FALCO_VERSION to match existing falco-driver-loader script.
    sed -i '/^\s*FALCO_VERSION=/d' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i "2 i FALCO_VERSION=\"${FALCO_VERSION}\"" "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i '3 i echo "FALCO_VERSION: $FALCO_VERSION"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Disable downloading from Falco driver repository.
    sed -i '/^\s*ENABLE_DOWNLOAD=/d' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i '2 i ENABLE_DOWNLOAD=""' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Enable compilation of eBPF driver.
    sed -i 's/ENABLE_COMPILE=.*/ENABLE_COMPILE="yes"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/DRIVER=.*/DRIVER="bpf"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Echo the KERNEL_RELEASE, KERNEL_VERSION and ARCH
    sed -i -e '/^KERNEL_RELEASE=.*/a\' -e 'echo "KERNEL_RELEASE: $KERNEL_RELEASE"' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i -e '/^KERNEL_VERSION=.*/a\' -e 'echo "KERNEL_VERSION: $KERNEL_VERSION"' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i -e '/^ARCH=.*/a\' -e 'echo "ARCH: $ARCH"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Set the KERNELDIR from the KERNEL_RELEASE. This is used in kernel Makefiles.
    sed -i -e '/^KERNEL_RELEASE=.*/a\' -e 'export KERNELDIR="/host/usr/src/kernels/$KERNEL_RELEASE"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Allow setting the outputs of `uname` via UNAME_* environment variables.
    sed -i 's/uname -r/echo "${UNAME_R}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -v/echo "${UNAME_V}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -m/echo "${UNAME_M}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    echo "Done!"

# Build falco probes in a recent version of Ubuntu to ensure we have up-to-date tooling
# containing required symbols (e.g. GLIBC_<RECENT_VERSION>).
FROM docker.io/ubuntu:${UBUNTU_VERSION}

ENV FALCO_DRIVER_LOADER_PATH="/usr/bin/falco-driver-loader"

# Set up the Falco enviromental variables.
ENV HOST_ROOT /host
ENV HOME /root

SHELL ["/bin/bash", "-c"]

# Install dev tools.
RUN set -Eeuxo pipefail; \
    apt-get update && apt-get install -y \
      build-essential \
      clang \
      curl \
      git \
      libelf-dev \
      llvm \
    && \
    rm -rf /var/lib/apt/lists/*

# Set up the Falco symlinks.
RUN rm -df /lib/modules \
	&& ln -s $HOST_ROOT/lib/modules /lib/modules

# Copy in entrypoint, falco source and falco-driver-loader script
COPY --from=falco-driver-loader "/docker-entrypoint.sh" "/docker-entrypoint.sh"
COPY --from=falco-driver-loader "/usr/src" "/usr/src"
COPY --from=falco-driver-loader "${FALCO_DRIVER_LOADER_PATH}" "${FALCO_DRIVER_LOADER_PATH}"

ENTRYPOINT ["/docker-entrypoint.sh"]
