go_library(
    name = "releasenotes",
    srcs = ["releasenotes.go"],
    visibility = [
        "//cmd/...",
    ],
)

go_test(
    name = "releasenotes_test",
    srcs = ["releasenotes_test.go"],
    external = True,
    deps = [
        ":releasenotes",
        "//third_party/go:stretchr_testify",
    ],
)
