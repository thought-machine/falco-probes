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
        "falcodriverbuilder_test.go",
    ],
    external = True,
    deps = [
        ":falcodriverbuilder",
        "//pkg/docker",
        "//third_party/go:stretchr_testify",
    ],
)
