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

go_module(
    name = "stretchr_testify",
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
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.0",
)

go_module(
    name = "pmezard_difflib",
    install = ["difflib"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_module(
    name = "stretchr_objx",
    licences = ["MIT"],
    module = "github.com/stretchr/objx",
    version = "v0.1.0",
    deps = [":davecgh_spew"],
)

go_module(
    name = "yaml_v3",
    licences = [
        "Apache-2.0",
        "MIT",
    ],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20200313102051-9f266ea9e77c",
)

go_module(
    name = "docker",
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
    licences = ["Apache-2.0"],
    module = "github.com/docker/go-units",
    version = "v0.4.0",
)

go_module(
    name = "docker_distribution",
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
    install = [
        "errdefs",
        "platforms",
        "log",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/containerd/containerd",
    version = "v1.5.2",
    deps = [
        ":google_grpc",
        ":opencontainers_imagespec",
        ":pkg_errors",
        ":sirupsen_logrus",
    ],
)

go_module(
    name = "pkg_errors",
    licences = ["BSD-2-Clause"],
    module = "github.com/pkg/errors",
    version = "v0.9.1",
)

go_module(
    name = "google_grpc",
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
    licences = ["MIT"],
    module = "github.com/sirupsen/logrus",
    version = "v1.8.1",
    deps = [
        ":x_sys",
    ],
)

go_module(
    name = "opencontainers_imagespec",
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
    licences = [
        "Apache-2.0",
        "CC-BY-SA-4.0",
    ],
    module = "github.com/opencontainers/go-digest",
    version = "v1.0.0",
)

go_module(
    name = "gogo_protobuf",
    install = [
        "proto",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/gogo/protobuf",
    version = "v1.3.2",
)

go_module(
    name = "x_net",
    install = [
        "proxy",
        "internal/socks",
        "http2",
        "http2/hpack",
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
    licences = ["BSD-3-Clause"],
    module = "github.com/jessevdk/go-flags",
    version = "v1.5.0",
    deps = [":x_sys"],
)

go_module(
    name = "moby_term",
    licences = ["Apache-2.0"],
    module = "github.com/moby/term",
    version = "v0.0.0-20210619224110-3f7ff695adc6",
    deps = [":x_sys"],
)

go_module(
    name = "morikuni_aec",
    licences = ["MIT"],
    module = "github.com/morikuni/aec",
    version = "v1.0.0",
)

go_module(
    name = "rs_zerolog",
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
    install = [
        "openpgp/...",
        "cast5",
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
