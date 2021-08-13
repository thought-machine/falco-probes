package releasenotes

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"text/template"

	"github.com/google/go-github/v37/github"
)

type ReleaseEditor interface {
	EditReleaseNotesByReleaseID(ctx context.Context, releaseID int64, body string) error
}

func EditReleaseNotes(ctx context.Context, releases []*github.RepositoryRelease, re ReleaseEditor) error {
	for _, r := range releases {
		probes := make(ReleasedProbes, len(r.Assets))
		for i, a := range r.Assets {
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
