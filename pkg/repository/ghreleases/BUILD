go_library(
    name = "ghreleases",
    srcs = ["ghreleases.go"],
    visibility = ["//build/..."],
    deps = [
        "//internal/logging",
        "//pkg/repository",
        "//third_party/go:google_github",
        "//third_party/go:x_oauth2",
    ],
)

go_test(
    name = "ghreleases_test",
    size = "large",
    srcs = ["ghreleases_test.go"],
    external = True,
    deps = [
        ":ghreleases",
        "//third_party/go:google_github",
        "//third_party/go:stretchr_testify",
        "//third_party/go:x_oauth2",
    ],
)
