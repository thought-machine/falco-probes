go_library(
    name = "docker",
    srcs = [
        "build.go",
        "docker.go",
        "image.go",
        "logs.go",
        "run.go",
        "volume.go",
    ],
    visibility = [
        "//cmd/...",
        "//pkg/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/operatingsystem",
        "//third_party/go:docker",
    ],
)

go_test(
    name = "docker_test",
    size = "large",
    srcs = [
        "docker_test.go",
    ],
    external = True,
    deps = [
        ":docker",
        "//pkg/operatingsystem",
        "//third_party/go:stretchr_testify",
    ],
)
