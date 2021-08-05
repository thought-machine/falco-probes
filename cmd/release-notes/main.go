package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/releasenotes"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"

	"github.com/google/go-github/v37/github"
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
		shortDriverVersion, err := getShortDriverVersionFromAsset(r)
		if err != nil {
			log.Fatal(err)
		}

		probes[i] = releasenotes.ReleasedProbe{
			DriverVersion: shortDriverVersion,
			Probe:         r.GetName(),
			KernelPackage: operatingsystem.KernelPackageFromProbeName(r.GetName()),
		}
	}

	sort.Sort(sort.Reverse(probes)) // most recently released drivers/kernel packages first

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

// getShortDriverVersionFromAsset uses the download url to extract the (shortened) driver version
func getShortDriverVersionFromAsset(a *github.ReleaseAsset) (string, error) {
	url := strings.TrimPrefix(a.GetBrowserDownloadURL(), "https://")
	urlSplit := strings.Split(url, "/")
	if len(urlSplit) < 6 || urlSplit[5] == "" {
		return "", fmt.Errorf("could not determine driver version from released asset at %s", url)
	}
	shortDriverVersion := urlSplit[5] // github.com/$org/$repo/releases/download/$releasename/$assetname

	return shortDriverVersion, nil
}
