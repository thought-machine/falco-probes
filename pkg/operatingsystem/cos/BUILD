go_library(
    name = "cos",
    srcs = [
        "git-reader.go",
        "kernel-package.go",
        "operating-system.go",
    ],
    visibility = [
        "//pkg/...",
    ],
    deps = [
        "//pkg/docker",
        "//pkg/operatingsystem",
        "//pkg/operatingsystem/cos/buildid",
        "//third_party/go:go_git",
    ],
)

go_test(
    name = "cos_test",
    size = "large",
    srcs = [
        "git-reader_test.go",
        "kernel-package_test.go",
        "operating-system_test.go",
    ],
    external = True,
    deps = [
        ":cos",
        "//pkg/docker",
        "//pkg/operatingsystem",
        "//pkg/operatingsystem/cos/mock",
        "//third_party/go:go_git",
        "//third_party/go:stretchr_testify",
    ],
)
