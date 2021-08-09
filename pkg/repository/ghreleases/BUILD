go_library(
    name = "ghreleases",
    srcs = [
        "cache.go",
        "ghreleases.go",
    ],
    visibility = [
        "//build/...",
        "//cmd/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/repository",
        "//third_party/go:google_github",
        "//third_party/go:x_oauth2",
    ],
)

go_test(
    name = "ghreleases_test",
    srcs = [
        "cache_test.go",
        "ghreleases_test.go",
    ],
    external = True,
    sandbox = False,
    deps = [
        ":ghreleases",
        "//third_party/go:google_github",
        "//third_party/go:stretchr_testify",
        "//third_party/go:x_oauth2",
    ],
)
