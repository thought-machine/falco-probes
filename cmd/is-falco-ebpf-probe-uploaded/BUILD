go_binary(
    name = "is-falco-ebpf-probe-uploaded",
    srcs = ["main.go"],
    deps = [
        "//internal/cmd",
        "//internal/logging",
        "//pkg/docker",
        "//pkg/falcodriverbuilder",
        "//pkg/operatingsystem/resolver",
        "//pkg/repository",
        "//pkg/repository/ghreleases",
    ],
)

# To Use:
#   Requires Personal access token with no scope selected to be added to .github_token
#   export GITHUB_TOKEN=$(cat .github_token)
#   plz run //cmd/is-falco-ebpf-probe-uploaded -- --falco_version <falco_version> <operating_system> <kernel_package_name>