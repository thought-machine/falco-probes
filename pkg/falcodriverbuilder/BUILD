go_library(
    name = "falcodriverbuilder",
    srcs = [
        "build-ebpf-probe.go",
        "falcodriverbuilder.go",
    ],
    resources = ["falco-driver-builder.Dockerfile"],
    visibility = [
        "//build/...",
        "//cmd/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/docker",
        "//pkg/operatingsystem",
    ],
)

go_test(
    name = "falcodriverbuilder_test",
    size = "large",
    srcs = [
        "build-ebpf-probe_test.go",
        "falcodriverbuilder_test.go",
    ],
    external = True,
    deps = [
        ":falcodriverbuilder",
        "//pkg/docker",
        "//pkg/operatingsystem/resolver",
        "//third_party/go:stretchr_testify",
    ],
)
