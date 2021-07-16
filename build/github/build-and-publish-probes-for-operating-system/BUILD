go_binary(
    name = "build-and-publish-probes-for-operating-system",
    srcs = ["main.go"],
    deps = [
        "//internal/cmd",
        "//internal/logging",
        "//pkg/docker",
        "//pkg/falcodriverbuilder",
        "//pkg/operatingsystem",
        "//pkg/operatingsystem/resolver",
        "//pkg/repository",
        "//pkg/repository/ghreleases",
    ],
)
