package main

import (
	"context"
	"time"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/releasenotes"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

var _ releasenotes.ReleaseEditor = (*ghreleases.GHReleases)(nil)

var log = logging.Logger

func main() {
	opts := &ghreleases.Opts{}
	cmd.MustParseFlags(opts)
	ghReleases := ghreleases.MustGHReleases(opts)

	releases, err := ghReleases.GetReleases()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get list of previously releases")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	if err := releasenotes.EditReleaseNotes(ctx, releases, ghReleases); err != nil {
		log.Fatal().Err(err).Msg("could not update release notes")
	}
}
