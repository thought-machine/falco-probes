sh_cmd(
    name = "go",
    srcs = ["go.sh"],
    cmd = """
source $(out_location //scripts/util)
export GOROOT=$(out_location //third_party/lang:go_tool)
export GO_FMT="\\\\$GOROOT/bin/gofmt"
source $SRCS
    """,
    data = [
        "//scripts/util",
        "//third_party/lang:go_tool",
    ],
    labels = ["fmt"],
    shell = "/bin/bash",
)

sh_cmd(
    name = "plz",
    srcs = ["plz.sh"],
    cmd = """
source $(out_location //scripts/util)
source $SRCS
    """,
    data = [
        "//scripts/util",
    ],
    labels = ["fmt"],
    shell = "/bin/bash",
)
