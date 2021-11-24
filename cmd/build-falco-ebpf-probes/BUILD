go_binary(
    name = "build-falco-ebpf-probes",
    srcs = ["main.go"],
    deps = [
        "//internal/cmd",
        "//internal/logging",
        "//internal/queue",
        "//pkg/docker",
        "//pkg/falcodriverbuilder",
        "//pkg/operatingsystem",
        "//pkg/operatingsystem/resolver",
    ],
)
