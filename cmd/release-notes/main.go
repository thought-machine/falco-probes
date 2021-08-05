package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/releasenotes"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

func main() {
	opts := &ghreleases.Opts{}
	cmd.MustParseFlags(opts)
	ghReleases := ghreleases.MustGHReleases(opts)

	released, err := ghReleases.GetReleasedProbes()
	if err != nil {
		log.Fatalf("unable to get list of previously released probes. %v", err)
	}

	probes := make(releasenotes.ReleasedProbes, len(released))
	for i, r := range released {
		url := strings.TrimPrefix(r.GetBrowserDownloadURL(), "https://")
		urlSplit := strings.Split(url, "/")
		if len(urlSplit) < 6 || urlSplit[5] == "" {
			log.Fatalf("could not determine driver version from released asset %s", url)
		}

		probes[i] = releasenotes.ReleasedProbe{
			DriverVersion: urlSplit[5],
			Probe:         r.GetName(),
			KernelPackage: operatingsystem.KernelPackageFromProbeName(r.GetName()),
		}
	}

	// Ensure we're grouping probes by driver version and then ordering by name
	sort.Sort(probes)

	tmpl, err := template.New("releasenotes").Parse(releaseNotesTemplate)
	if err != nil {
		log.Fatalf("unable to parse release notes template. %v", err)
	}

	params := releaseNoteParams{ProbeRows: make([]string, len(probes))}
	for i, p := range probes {
		params.ProbeRows[i] = p.ToMarkdownRow()
	}

	w, err := os.Create("RELEASE_NOTES.md")
	if err != nil {
		log.Fatalf("unable to create RELEASE_NOTES.md. %v", err)
	}
	defer w.Close()

	if err := tmpl.Execute(w, &params); err != nil {
		log.Fatalf("unable to write probes to RELEASE_NOTES.md. %v", err)
	}
}
