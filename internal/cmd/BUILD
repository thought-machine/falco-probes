go_library(
    name = "cmd",
    srcs = [
        "flags.go",
        "logging.go",
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
