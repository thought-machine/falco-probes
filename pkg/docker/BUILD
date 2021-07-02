go_library(
    name = "docker",
    srcs = [
        "docker.go",
        "logs.go",
        "run.go",
    ],
    visibility = ["//pkg/..."],
    deps = [
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
        "//third_party/go:stretchr_testify",
    ],
)
