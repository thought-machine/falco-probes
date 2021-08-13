go_library(
    name = "releasenotes",
    srcs = [
        "releasenotes.go",
        "releasedprobe.go",
        "tmpl.go",
    ],
    deps = [
        "//third_party/go:google_github",
    ],
    visibility = ["//cmd/..."],
)

go_test(
    name = "releasenotes_test",
    srcs = [
        "releasenotes_test.go",
        "releasedprobe_test.go",
    ],
    external = True,
    deps = [
        ":releasenotes",
        "//third_party/go:stretchr_testify",
    ],
)