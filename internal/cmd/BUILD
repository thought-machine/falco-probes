go_library(
    name = "cmd",
    srcs = [
        "flags.go",
        "logging.go",
        "parallel.go",
    ],
    visibility = [
        "//build/...",
        "//cmd/...",
    ],
    deps = [
        "//internal/logging",
        "//third_party/go:jessevdk_flags",
        "//third_party/go:rs_zerolog",
    ],
)

go_test(
    name = "cmd_test",
    srcs = [
        "parallel_test.go",
    ],
    external = True,
    deps = [
        ":cmd",
        "//third_party/go:stretchr_testify",
    ],
)
