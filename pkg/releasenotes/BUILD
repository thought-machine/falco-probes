go_library(
    name = "releasenotes",
    srcs = [
        "releasedprobe.go",
        "releasenotes.go",
        "tmpl.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "//third_party/go:google_github",
    ],
)

go_test(
    name = "releasenotes_test",
    srcs = [
        "releasedprobe_test.go",
        "releasenotes_test.go",
    ],
    external = True,
    deps = [
        ":releasenotes",
        "//third_party/go:google_github",
        "//third_party/go:stretchr_testify",
    ],
)
