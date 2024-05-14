go_library(
    name = "operatingsystem",
    srcs = [
        "kernel-package.go",
        "operating-system.go",
    ],
    deps = [
        "//internal/logging",
    ],
    visibility = [
        "//build/...",
        "//cmd/...",
        "//pkg/...",
    ],
)

go_test(
    name = "operatingsystem_test",
    srcs = ["kernel-package_test.go"],
    external = True,
    deps = [
        ":operatingsystem",
        "//third_party/go:stretchr_testify",
    ],
)
