go_library(
    name = "queue",
    srcs = [
        "queue.go",
        "worker.go",
        "task.go",
        "opts.go",
    ],
    visibility = [
        "//cmd/...",
    ],
    deps = [
        "//internal/logging",
    ],
)

go_test(
    name = "queue_test",
    srcs = [
        "queue_test.go",
    ],
    external = True,
    deps = [
        ":queue",
        "//third_party/go:stretchr_testify",
    ],
)
