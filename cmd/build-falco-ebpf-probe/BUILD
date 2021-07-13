go_binary(
    name = "build-falco-ebpf-probe",
    srcs = ["main.go"],
    deps = [
        "//internal/cmd",
        "//internal/logging",
        "//pkg/docker",
        "//pkg/falcodriverbuilder",
        "//pkg/operatingsystem",
        "//pkg/operatingsystem/resolver",
    ],
)
