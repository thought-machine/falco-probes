go_library(
    name = "amazonlinux2",
    srcs = [
        "kernel-package.go",
        "operating-system.go",
    ],
    deps = [
        "//pkg/docker",
        "//pkg/operatingsystem",
    ],
)

go_test(
    name = "amazonlinux2_test",
    size = "large",
    srcs = [
        "kernel-package_test.go",
        "operating-system_test.go",
    ],
    external = True,
    deps = [
        ":amazonlinux2",
        "//pkg/docker",
        "//pkg/operatingsystem",
        "//third_party/go:stretchr_testify",
    ],
)
