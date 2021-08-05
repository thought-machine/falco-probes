package main

type releaseNoteParams struct {
	ProbeRows []string
}

var releaseNotesTemplate = `
# Released Probes
| Kernel Package | Probe |
|----------------|-------|
{{range .ProbeRows}}{{.}}
{{end}}
`
