package releasenotes

type releaseNotesParams struct {
	ProbeRows []string
}

var releaseNotesTemplate = `
# Probes
| Kernel Package | Probe |
|----------------|-------|
{{range .ProbeRows}}{{.}}
{{end}}`
