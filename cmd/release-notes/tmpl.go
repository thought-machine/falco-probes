package main

type releaseNoteParams struct {
	ProbeRows []string
}

var releaseNotesTemplate = `
# Released Probes
| Driver | Kernel Package | Probe |
|--------|----------------|-------|
{{range .ProbeRows}}{{.}}
{{end}}
`
