package main

type releasedProbesParams struct {
	ProbeRows []string
}

var releasedProbesTemplate = `
# Released Probes
| Kernel Package | Probe |
|----------------|-------|
{{range .ProbeRows}}{{.}}
{{end}}
`
