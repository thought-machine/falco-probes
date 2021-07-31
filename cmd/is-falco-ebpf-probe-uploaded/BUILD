go_binary(
    name = "is-falco-ebpf-probe-uploaded",
    srcs = ["main.go"],
    deps = [
        "//internal/cmd",
        "//internal/logging",
        "//pkg/docker",
        "//pkg/falcodriverbuilder",
        "//pkg/repository",
        "//pkg/repository/ghreleases",
        "//pkg/operatingsystem/resolver",
    ],
)

# E.g. plz run //cmd/is-falco-ebpf-probe-uploaded -- --verbose --falco_version=0.29.1 --github_releases_token=$GRT amazonlinux2 4.14.200-155.322.amzn2
