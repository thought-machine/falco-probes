package releasenotes

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/google/go-github/v37/github"
)

// ReleaseEditor abstracts the implementation of editing an existing github release
type ReleaseEditor interface {
	EditReleaseNotesByReleaseID(ctx context.Context, releaseID int64, body string) error
}

// ParseProbesFromReleaseNotes converts the markdown table from the body of the release into a list of probes
func ParseProbesFromReleaseNotes(release *github.RepositoryRelease) ReleasedProbes {
	body := release.GetBody()
	probes := ReleasedProbes{}

	for _, line := range strings.Split(body, "\n") {
		p := ReleasedProbeFromMarkdownRow(line)

		if p.KernelPackage != "" {
			probes = append(probes, p)
		}
	}

	return probes
}

// ListKernelPackagesToCompile cross-references the list of released probes against the provided kernel
// package names, returning a list of all the kernels that have not been released for _all_ Falco versions.
func (rps ReleasedProbes) ListKernelPackagesToCompile(kernelPackageNames []string, numReleases int) []string {
	toCompile := kernelPackageNames[:0]
	for _, kpn := range kernelPackageNames {
		numProbesMissing := numReleases

		for _, rp := range rps {
			if rp.KernelPackage == kpn {
				numProbesMissing--
			}
		}

		if numProbesMissing > 0 {
			toCompile = append(toCompile, kpn)
		}
	}

	return toCompile
}

// SetReleaseNotes updates the provided releases, using custom templating logic, via the github API
func SetReleaseNotes(ctx context.Context, releases []*github.RepositoryRelease, re ReleaseEditor) error {
	for _, r := range releases {
		probes := make(ReleasedProbes, len(r.Assets))
		for i, a := range r.Assets {
			if !strings.HasSuffix(a.GetName(), ".o") {
				continue
			}

			probes[i] = ReleasedProbe{
				Probe:         a.GetName(),
				KernelPackage: KernelPackageFromProbeName(a.GetName()),
			}
		}

		sort.Sort(sort.Reverse(&probes)) // Sort our probes by most recent kernel version first

		params := releaseNotesParams{ProbeRows: make([]string, len(probes))}
		for i, p := range probes {
			params.ProbeRows[i] = p.ToMarkdownRow()
		}

		tmpl, err := template.New("releasenotes").Parse(releaseNotesTemplate)
		if err != nil {
			return fmt.Errorf("unable to parse release notes template. %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, &params); err != nil {
			return fmt.Errorf("unable to template release note content. %w", err)
		}

		if err := re.EditReleaseNotesByReleaseID(ctx, r.GetID(), buf.String()); err != nil {
			return fmt.Errorf("error updating release notes with list of released probes. %w", err)
		}
	}

	return nil
}
