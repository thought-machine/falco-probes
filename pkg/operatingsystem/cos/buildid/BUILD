go_library(
    name = "buildid",
    srcs = [
        "validator.go",
    ],
    visibility = [
        "//pkg/...",
    ],
)

go_test(
    name = "buildid_test",
    size = "large",
    srcs = [
        "validator_test.go",
    ],
    external = True,
    deps = [
        "//pkg/operatingsystem/cos/buildid",
        "//pkg/operatingsystem/cos/mock",
        "//third_party/go:stretchr_testify",
    ],
)
