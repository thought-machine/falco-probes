ARG GENTOO_VERSION=20221031

# Add portage with pre-populated gentoo repo.
FROM "gentoo/portage:${GENTOO_VERSION}" as portage

FROM "gentoo/stage3:${GENTOO_VERSION}"

ARG COS_REPO_URLS="chromiumos-overlay eclass-overlay portage-stable"
ARG COS_REPO_NAMES="chromiumos eclass-overlay portage-stable"
ENV COS_REPO_URLS="${COS_REPO_URLS}"
ENV COS_REPO_NAMES="${COS_REPO_NAMES}"

# Copy the entire portage gentoo repo in.
COPY --from=portage /var/db/repos/gentoo /var/db/repos/gentoo

RUN emerge -qv app-misc/jq dev-vcs/git

RUN set -Eeuxo pipefail; \
    # Set COS repos.
    REPO_URLS=(${COS_REPO_URLS}) && \
    REPO_NAMES=(${COS_REPO_NAMES}) && \
    # Add COS repos.
    mkdir -p /etc/portage/repos.conf && \
    for ((i=0; i<${#REPO_URLS[@]}; i++)); do \
        repo_url=${REPO_URLS[$i]} && \
        repo_name=${REPO_NAMES[$i]} && \
        repo_path=/var/db/repos/${repo_name} && \
        git clone https://cos.googlesource.com/third_party/overlays/${repo_url} /var/db/repos/${repo_name} && \
        printf '[%s]\nlocation = %s\npriority = 100\n' ${repo_name} ${repo_path} \
          > /etc/portage/repos.conf/${repo_name}-repo.conf; done && \
    echo "Done!" \
