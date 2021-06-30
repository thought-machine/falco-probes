go_module(
    name = "x_tools",
    hashes = ["dff907b66b1f2cfd3536c92d71dfd5f28078026249f850214ee1a860519decc6"],
    install = [
        "cover",
        "present",
        "imports",
        "go/packages",
        "go/internal/packagesdriver",
        "internal/gocommand",
        "internal/packagesinternal",
        "go/gcexportdata",
        "go/internal/gcimporter",
        "go/ast/astutil",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/tools",
    version = "5bcca83a78812bd91ce8cb29be4fc4521cdc8f6f",
    visibility = ["PUBLIC"],
    deps = [
        ":x_errors",
        ":x_mod",
    ],
)

go_module(
    name = "x_errors",
    hashes = ["10780b586a60adfdea2f4d76dbaf3025eb3a8165376153c7febfb36516d1ec72"],
    install = [
        ".",
        "...",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/xerrors",
    version = "a5947ffaace3e882f334c1750858b4a6a7e52422",
    deps = [":x_sys"],
)

go_module(
    name = "x_sys",
    hashes = ["e9ba045ac49a174a9bb1220fc50e89b46e26c8ea18d987248baec06478c9c7bf"],
    install = [
        "unix",
        "cpu",
    ],
    module = "golang.org/x/sys",
    version = "95b1ffbd15a57cc5abb3f04402b9e8ec0016a52c",
)

go_module(
    name = "x_mod",
    hashes = ["2a937734dda11bfc94447ce43128e5ffafcf0479a4f4a4558aa17baf7335c2d7"],
    install = [
        "semver",
        "module",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/mod",
    version = "v0.4.2",
    deps = [":x_errors"],
)
