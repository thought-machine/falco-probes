package main

import (
	"log"
	"text/template"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/pkg/releasenotes"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

func main() {
	opts := &ghreleases.Opts{}
	cmd.MustParseFlags(opts)
	ghReleases := ghreleases.MustGHReleases(opts)

	releases, err := ghReleases.GetReleases()
	if err != nil {
		log.Fatalf("unable to get list of previously releases. %v", err)
	}

	probes := make(releasenotes.ReleasedProbes, 0)
	for _, r := range releases {
		for _, a := range r.Assets {
			probes = append(probes, releasenotes.ReleasedProbe{
				Probe:         a.GetName(),
				KernelPackage: releasenotes.KernelPackageFromProbeName(a.GetName()),
			})
		}
	}

	params := releaseNoteParams{ProbeRows: make([]string, len(probes))}
	for i, p := range probes {
		params.ProbeRows[i] = p.ToMarkdownRow()
	}

	tmpl, err := template.New("releasenotes").Parse(releaseNotesTemplate)
	if err != nil {
		log.Fatalf("unable to parse release notes template. %v", err)
	}

}
