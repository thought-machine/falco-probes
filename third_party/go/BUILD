package(default_visibility = ["PUBLIC"])

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
    hashes = ["ef9844802ee107d2021908789cfdbc0f321972ac8cc0e4232c4f204fdd40434e"],
    install = [
        "cpu",
        "execabs",
        "internal/unsafeheader",
        "unix",
    ],
    module = "golang.org/x/sys",
    version = "v0.0.0-20210630005230-0f9fa26af87c",
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

go_module(
    name = "stretchr_testify",
    hashes = ["bf8d2167ee8dbc994d482a0fa51f0e8ef341a2d3f98e9370b3deaed105e2c2e4"],
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.7.0",
    deps = [
        ":davecgh_spew",
        ":pmezard_difflib",
        ":stretchr_objx",
        ":yaml_v3",
    ],
)

go_module(
    name = "davecgh_spew",
    hashes = ["221a963c9be5f3352497b7310c43435df7ae3827d03e4c948cda3bce46f8a5aa"],
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.0",
)

go_module(
    name = "pmezard_difflib",
    hashes = ["390ffaa52adc6904c9a963e7aab1aaaf5e2127088f502fbdee297fad201264af"],
    install = ["difflib"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_module(
    name = "stretchr_objx",
    hashes = ["d4cc594e30d93bfb6a92fe69aeb39f43cfba21fc76ba472d5b2b8bf857200691"],
    licences = ["MIT"],
    module = "github.com/stretchr/objx",
    version = "v0.1.0",
    deps = [":davecgh_spew"],
)

go_module(
    name = "yaml_v3",
    hashes = ["93964f906a898a9446b12139755d894dca4ded23428f81d2da35e4d119feb9b6"],
    licences = [
        "Apache-2.0",
        "MIT",
    ],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20200313102051-9f266ea9e77c",
)

go_module(
    name = "docker",
    hashes = ["9c5f69b2668463ec5d448cf2326c8099ce4b329901e914660564e909ea4dade4"],
    install = [
        "api/types",
        "api/types/container",
        "api/types/filters",
        "api/types/versions",
        "api/types/swarm",
        "api/types/swarm/runtime",
        "api/types/network",
        "api/types/time",
        "api/types/registry",
        "api/types/events",
        "api/types/image",
        "api/types/volume",
        "api/types/mount",
        "api/types/strslice",
        "api/types/blkiodev",
        "api",
        "client",
        "pkg/stdcopy",
        "pkg/jsonmessage",
        "errdefs",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/docker/docker",
    version = "v20.10.7+incompatible",
    deps = [
        ":containerd",
        ":docker_connections",
        ":docker_distribution",
        ":docker_units",
        ":gogo_protobuf",
        ":moby_term",
        ":morikuni_aec",
        ":opencontainers_imagespec",
        ":sirupsen_logrus",
    ],
)

go_module(
    name = "docker_connections",
    hashes = ["06ea9101e470cd34d88bd56624547797456ad1c183b7861ea9f2df4b309d6df6"],
    install = [
        "nat",
        "sockets",
        "tlsconfig",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/docker/go-connections",
    version = "v0.4.0",
    deps = [
        ":pkg_errors",
        ":x_net",
    ],
)

go_module(
    name = "docker_units",
    hashes = ["8184490955536f93ec58b92f7b7e94dec9bb47c1fb9e7b9ae3825982a5618483"],
    licences = ["Apache-2.0"],
    module = "github.com/docker/go-units",
    version = "v0.4.0",
)

go_module(
    name = "docker_distribution",
    hashes = ["7ff1039ca0621f1a0a26c852bbfb969be8ead0860b65e61e20d5489f7cde390a"],
    install = [
        "registry/api/errcode",
        "reference",
        "digestset",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/docker/distribution",
    version = "v2.7.1+incompatible",
    deps = [
        ":opencontainers_digest",
    ],
)

go_module(
    name = "containerd",
    hashes = ["2b199a3f1b32543a8e02678c225307e457ca6d7ac3cb4d13545dea1d7614a81e"],
    install = [
        "errdefs",
        "platforms",
        "log",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/containerd/containerd",
    version = "v1.5.4",
    deps = [
        ":google_grpc",
        ":opencontainers_imagespec",
        ":pkg_errors",
        ":sirupsen_logrus",
    ],
)

go_module(
    name = "pkg_errors",
    hashes = ["2d51fe9d61d9879d44cf034c4b63b20bd967b68132e70490618950de478a56b2"],
    licences = ["BSD-2-Clause"],
    module = "github.com/pkg/errors",
    version = "v0.9.1",
)

go_module(
    name = "google_grpc",
    hashes = ["bf11c16decee5c08b5382eb4c8e99b4b19b90936f57a56a9d9d1fab22f8a9f05"],
    install = [
        "codes",
        "status",
        "internal/...",
    ],
    licences = ["Apache-2.0"],
    module = "google.golang.org/grpc",
    version = "v1.33.2",
    deps = [
        ":github_protobuf",
        ":google_genproto",
        ":x_net",
        ":x_sys",
    ],
)

go_module(
    name = "github_protobuf",
    hashes = ["6aaf398e1e116548f81bb4f1f94e9d8c59cfcc3179c09791847eeed936576a49"],
    install = [
        "proto",
        "ptypes",
        "ptypes/any",
        "ptypes/duration",
        "ptypes/timestamp",
        "ptypes/wrappers",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/golang/protobuf",
    version = "v1.5.2",
    deps = [
        ":google_protobuf",
    ],
)

go_module(
    name = "google_protobuf",
    hashes = ["5ed9c2a320ca4ec4d7bf2f0d06d8da760c8a2849929cf10f9e7353bfdb2c7511"],
    install = [
        "proto",
        "encoding/prototext",
        "encoding/protowire",
        "runtime/protoimpl",
        "runtime/protoiface",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "reflect/protodesc",
        "types/known/anypb",
        "types/known/durationpb",
        "types/known/timestamppb",
        "types/known/wrapperspb",
        "types/descriptorpb",
        "internal/...",
    ],
    licences = ["BSD-3-Clause"],
    module = "google.golang.org/protobuf",
    version = "v1.27.1",
)

go_module(
    name = "google_genproto",
    hashes = ["a0288e8e857babe7417f595e31c41eba207d32f7637feea72924eb8eb8f9c24c"],
    install = [
        "googleapis/rpc/status",
        "googleapis/rpc/code",
    ],
    licences = ["Apache-2.0"],
    module = "google.golang.org/genproto",
    version = "v0.0.0-20200526211855-cb27e3aa2013",
    deps = [
        ":github_protobuf",
    ],
)

go_module(
    name = "sirupsen_logrus",
    hashes = ["cabe686e3c65689211cec3049da00ec84dc63e1d8ad93078c4951ff2a2337e4c"],
    licences = ["MIT"],
    module = "github.com/sirupsen/logrus",
    version = "v1.8.1",
    deps = [
        ":x_sys",
    ],
)

go_module(
    name = "opencontainers_imagespec",
    hashes = ["0978c46a1e50b95fd5b1c47cc5a6b08a5ae20ca6d4dddb395e6c1c56344bcfab"],
    install = [
        "specs-go/v1",
        "specs-go",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/opencontainers/image-spec",
    version = "v1.0.1",
    deps = [
        ":opencontainers_digest",
    ],
)

go_module(
    name = "opencontainers_digest",
    hashes = ["b01b5756818ec5678b7eb01ac198b0cfeb13ded384427e62d88f46b07f929a37"],
    licences = [
        "Apache-2.0",
        "CC-BY-SA-4.0",
    ],
    module = "github.com/opencontainers/go-digest",
    version = "v1.0.0",
)

go_module(
    name = "gogo_protobuf",
    hashes = ["14b27e769cbfcef5e9311461ff3767b91a6766c697b7a847c250e01ee00cfb52"],
    install = [
        "proto",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/gogo/protobuf",
    version = "v1.3.2",
)

go_module(
    name = "x_net",
    hashes = ["dea86856da5654da44af6c5111a75d2e4e43e6476601f94be2a2534b4a19721c"],
    install = [
        "proxy",
        "internal/socks",
        "http2",
        "http2/hpack",
        "context",
        "context/ctxhttp",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/net",
    version = "v0.0.0-20210614182718-04defd469f4e",
    deps = [
        ":x_text",
    ],
)

go_module(
    name = "x_text",
    hashes = ["88426eed435529e8f1953181958383d891ff6e70b51e3143b51ef25b76a6b33f"],
    install = [
        "secure/bidirule",
        "unicode/bidi",
        "unicode/norm",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/text",
    version = "v0.3.6",
)

go_module(
    name = "jessevdk_flags",
    hashes = ["0360c29c8a9085e6e60c30c5057c694509a523136f7c95dfbdf4434fca79c9b3"],
    licences = ["BSD-3-Clause"],
    module = "github.com/jessevdk/go-flags",
    version = "v1.5.0",
    deps = [":x_sys"],
)

go_module(
    name = "moby_term",
    hashes = ["996e71c68a725c6e30053d171595e74e3e75b724332fe1f3cd759655a1bbdb49"],
    licences = ["Apache-2.0"],
    module = "github.com/moby/term",
    version = "v0.0.0-20210619224110-3f7ff695adc6",
    deps = [":x_sys"],
)

go_module(
    name = "morikuni_aec",
    hashes = ["aabeaf991b8e5d32c048301c67ba440cd0e7194dd35dc1a2634eb056e1c65aa5"],
    licences = ["MIT"],
    module = "github.com/morikuni/aec",
    version = "v1.0.0",
)

go_module(
    name = "rs_zerolog",
    hashes = ["768e9946fb4186191e0de1563987f304d8483764dd633a79e259129b2564f0a2"],
    install = [
        ".",
        "internal/json",
    ],
    licences = ["MIT"],
    module = "github.com/rs/zerolog",
    version = "v1.23.0",
)

go_module(
    name = "google_github",
    hashes = ["071671143e2aee8c5d2df516de3d2516961c1125e32bc8cb5958359bc819bbcb"],
    install = ["github"],
    licences = ["BSD-3-Clause"],
    module = "github.com/google/go-github/v37",
    version = "v37.0.0",
    deps = [
        ":google_querystring",
        ":x_crypto",
    ],
)

go_module(
    name = "x_crypto",
    hashes = ["6de6deee143698ea03008825cd80eea594867f70c7ae435609d01df420f9b3a7"],
    install = [
        "blowfish",
        "cast5",
        "chacha20",
        "curve25519",
        "curve25519/...",
        "ed25519",
        "internal/...",
        "openpgp/...",
        "poly1305",
        "ssh",
        "ssh/agent",
        "ssh/knownhosts",
        "ssh/internal/...",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/crypto",
    version = "v0.0.0-20210711020723-a769d52b0f97",
)

go_module(
    name = "google_querystring",
    install = ["query"],
    licences = ["BSD-3-Clause"],
    module = "github.com/google/go-querystring",
    version = "v1.1.0",
)

go_module(
    name = "x_oauth2",
    install = [
        ".",
        "internal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/oauth2",
    version = "v0.0.0-20210628180205-a41e5a781914",
    deps = [":x_net"],
)

go_module(
    name = "diff",
    hashes = ["f33c2c7e4e2d287b933a906e19b7a93a7c131864010aba0e5f5ae8e4843199d9"],
    licences = ["BSD-3-Clause"],
    module = "github.com/pkg/diff",
    version = "v0.0.0-20210226163009-20ebb0f2a09e",
)

go_module(
    name = "go_billy",
    hashes = ["6d6291ad22df48bc2c371f47c5bfde6cb8ea025d050b4fb7d80e6fb8527bd31a"],
    install = [
        ".",
        "helper/...",
        "osfs",
        "util",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/go-git/go-billy/v5",
    version = "v5.3.1",
    deps = [":x_sys"],
)

go_module(
    name = "go_check",
    hashes = ["d18350b1be1eaffdd684775792da48f87622d1231015aaae939d45329723cc66"],
    install = ["."],
    licences = ["BSD-2-Clause"],
    module = "gopkg.in/check.v1",
    version = "v1.0.0-20201130134442-10cb98267c6c",
    deps = [":pretty"],
)

go_module(
    name = "go_context",
    hashes = ["c50f1e227d5f1c6b56b35f7bbe379a2ba60eaacae52854fdc7a0e2cf860c8b84"],
    install = [
        ".",
        "io",
    ],
    licences = ["MIT"],
    module = "github.com/jbenet/go-context",
    version = "v0.0.0-20150711004518-d14ea06fba99",
    deps = [":x_net"],
)

go_module(
    name = "go_crypto",
    hashes = ["c9b13257caa8309c233bc6952dc42c3b0e5e906aa2b6bc16f2173001b6edb107"],
    install = [
        "bitcurves",
        "brainpool",
        "eax",
        "internal/...",
        "ocb",
        "openpgp",
        "openpgp/...",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/ProtonMail/go-crypto",
    version = "v0.0.0-20210428141323-04723f9f07d7",
    deps = [":x_crypto"],
)

go_module(
    name = "go_data_structures",
    hashes = ["ece160c77a77ffe23321eea210bbd294ad64bb088efdcb98b27ba3951bdbabcc"],
    install = [
        "containers",
        "containers/...",
        "lists",
        "lists/...",
        "trees",
        "trees/...",
        "utils",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/emirpasic/gods",
    version = "v1.12.0",
)

go_module(
    name = "go_diff",
    hashes = ["296d1f94c0e11d86d2366a3af9f49ea93f2b4e2c0ac5e95c86486b1f29553ad3"],
    install = ["diffmatchpatch"],
    licences = ["MIT"],
    module = "github.com/sergi/go-diff",
    version = "v1.1.0",
)

go_module(
    name = "go_homedir",
    hashes = ["9b701a164a63068b56ded9e635e8a3602bbbcf13e2c5bf4ad5cd4f21d60e64d1"],
    install = ["."],
    licences = ["MIT"],
    module = "github.com/mitchellh/go-homedir",
    version = "v1.1.0",
)

go_module(
    name = "go_git",
    hashes = ["1db17db8dc7e8a6c3c4585fced3da7604a2e7b07f18d088aaae9a87fb2cf5631"],
    install = [
        ".",
        "config",
        "internal/...",
        "plumbing",
        "plumbing/...",
        "storage/...",
        "utils/...",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/go-git/go-git/v5",
    version = "v5.4.2",
    # https://github.com/go-git/go-git/blob/v5.4.2/go.mod
    deps = [
        ":go_billy",
        ":go_context",
        ":go_crypto",
        ":go_data_structures",
        ":go_diff",
        ":go_git_fixtures",
        ":go_git_gcfg",
        ":go_homedir",
        ":mergo",
        ":ssh_agent",
        ":ssh_config",
        ":x_crypto",
        ":x_sys",
    ],
)

go_module(
    name = "go_git_fixtures",
    hashes = ["8c54a94000a6e569254bd9734e802a7c833312df9bf82de929962c078e304522"],
    install = ["."],
    licences = ["Apache-2.0"],
    module = "github.com/go-git/go-git-fixtures/v4",
    version = "v4.2.1",
    deps = [
        ":go_billy",
        ":go_check",
    ],
)

go_module(
    name = "go_git_gcfg",
    hashes = ["145b0274ce955a383b632403aab831134a0df9d0cf5c5a840dd52b09544cf49b"],
    install = [
        ".",
        "scanner",
        "token",
        "types",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/go-git/gcfg",
    version = "v1.5.0",
    deps = [":go_warnings"],
)

go_module(
    name = "go_internal",
    hashes = ["ceb57d4dec298371fd7c08d307f0a653af95c9c6ae4e4dcb88f93a76edcf987d"],
    install = ["dirhash"],
    licences = ["BSD-3-Clause"],
    module = "github.com/rogpeppe/go-internal",
    version = "v1.9.0",
    deps = [":diff"],
)

go_module(
    name = "go_warnings",
    hashes = ["bad7f24b625c153413217c5be0f21f9204def34e61442735eb7a4b2dce5befd7"],
    install = ["."],
    licences = [
        "BSD-3-Clause",
        "ISC",
    ],
    module = "gopkg.in/warnings.v0",
    version = "v0.1.2",
)

go_module(
    name = "mergo",
    hashes = ["cb908075f66d9e0ba9755cc91f22519c050f61ab039e8904250b30ed5d35a350"],
    install = ["."],
    licences = ["BSD-3-Clause"],
    module = "github.com/imdario/mergo",
    version = "v0.3.12",
)

go_module(
    name = "pretty",
    hashes = ["d4ec83c3003d40440883a6a2ebfa2acd5ce5e74f196350f6c73aff3baa5f76e2"],
    licences = ["MIT"],
    module = "github.com/kr/pretty",
    version = "v0.2.1",
    deps = [
        ":go_internal",
        ":text",
    ],
)

go_module(
    name = "pty",
    hashes = ["ccaba64845b9896b17dc0435cca82cec107a2166814df3856115c138fdf3a792"],
    licences = ["MIT"],
    module = "github.com/creack/pty",
    version = "v1.1.17",
)

go_module(
    name = "ssh_agent",
    hashes = ["ae2a4ee6e7e6b3fc8c8dc275c19adf14346ed540f609dfb5dcf24116c8e262e6"],
    install = ["."],
    licences = ["Apache-2.0"],
    module = "github.com/xanzy/ssh-agent",
    version = "v0.3.0",
    deps = [":x_crypto"],
)

go_module(
    name = "ssh_config",
    hashes = ["6b97c18442738d8ee56e03c0bf8035290bfaa8b6622956efe1b4373525c744bf"],
    install = ["."],
    licences = [
        "BSD-3-Clause",
        "MIT",
    ],
    module = "github.com/kevinburke/ssh_config",
    version = "v0.0.0-20201106050909-4977a11b4351",
)

go_module(
    name = "text",
    hashes = ["76c6a65454081b1b7d5996d576a67ab58b47f6b57ec843f1d811b5f0de4835cd"],
    install = ["."],
    licences = ["MIT"],
    module = "github.com/kr/text",
    version = "v0.2.0",
    deps = [
        ":pty",
    ],
)
