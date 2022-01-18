ARG FALCO_VERSION
FROM "docker.io/falcosecurity/falco-driver-loader:${FALCO_VERSION}"

ENV FALCO_DRIVER_LOADER_PATH="/usr/bin/falco-driver-loader"

SHELL ["/bin/bash", "-c"]

RUN set -Eeuxo pipefail; \
    # Determine DRIVERS_REPO, DRIVER_VERSION, DRIVER_NAME from existing falco-driver-loader to persist them in the 0.30.0 falco-driver-loader script.
    export DRIVERS_REPO=$(grep ^DRIVERS_REPO "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    export DRIVER_VERSION=$(grep ^DRIVER_VERSION "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    export DRIVER_NAME=$(grep ^DRIVER_NAME "${FALCO_DRIVER_LOADER_PATH}" | cut -f2 -d\") && \
    # Use falco-driver-loader from 0.30.0 (the patches below work with that script.)
    curl -L https://raw.githubusercontent.com/falcosecurity/falco/0.30.0/scripts/falco-driver-loader \
    -o /usr/bin/falco-driver-loader && chmod +x /usr/bin/falco-driver-loader && \
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
    # Set the KERNELDIR from the KERNEL_RELEASE. This is used in kernel Makefiles.
    sed -i -e '/^KERNEL_RELEASE=.*/a\' -e 'export KERNELDIR="/lib/modules/$KERNEL_RELEASE/build"' "${FALCO_DRIVER_LOADER_PATH}" && \
    # Allow setting the outputs of `uname` via UNAME_* environment variables.
    sed -i 's/uname -r/echo "${UNAME_R}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -v/echo "${UNAME_V}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    sed -i 's/uname -m/echo "${UNAME_M}"/g' "${FALCO_DRIVER_LOADER_PATH}" && \
    echo "Done!"
