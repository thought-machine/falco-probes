package main

import (
	"os"
	"strings"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
	"github.com/thought-machine/falco-probes/pkg/repository"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

// Takes same inputs as the //cmd/build-falco-ebpf-probe tool.
// TODO: figue out how to deal with the GHReleases issue & 401 for calling github api
type opts struct {
	FalcoVersion string          `long:"falco_version" description:"The version of Falco to compile probes against" required:"true"`
	GHReleases   ghreleases.Opts `group:"github_releases" namespace:"github_releases"`
	Positional   struct {
		OperatingSystem string `positional-arg-name:"operating_system"`
		KernelPackage   string `positional-arg-name:"kernel_package"`
	} `positional-args:"yes" required:"true"`
}

var log = logging.Logger

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)

	cli := docker.MustClient()
	var repo repository.Repository = ghreleases.MustGHReleases(&opts.GHReleases)

	// Verify the inputs
	log.Info().Str("falco_version", opts.FalcoVersion).Msg("Verifying input")
	falcoDriverBuilderImg, err := falcodriverbuilder.BuildImage(cli, opts.FalcoVersion)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get driver builder image for provided falco_version")
	}
	driverVersion, err := falcodriverbuilder.GetDriverVersion(cli, falcoDriverBuilderImg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get driver version for provided falco_version")
	}

	log.Info().Str("operating_system", opts.Positional.OperatingSystem).Msg("Verifying input")
	operatingSystem, err := resolver.OperatingSystem(cli, opts.Positional.OperatingSystem)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get operating system")
	}

	log.Info().Str("kernel_package", opts.Positional.KernelPackage).Msg("Verifying input")
	kernelPackage, err := operatingSystem.GetKernelPackageByName(opts.Positional.KernelPackage)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get kernel package")
	}
	probeName := kernelPackage.ProbeName() + ".o"
	if strings.Contains(probeName, operatingSystem.GetName()) == false {
		log.Fatal().Msg("Inputted probe doesn't match inputted operating system.")
	}

	// Indentify if probe uploaded
	log.Info().
		Str("driverVersion", driverVersion).
		Str("probeName", probeName).
		Msg("Identifying if falco ebpf probe uploaded for ")

	isAlreadyMirrored, err := repo.IsAlreadyMirrored(driverVersion, probeName)
	if err != nil {
		log.Fatal().Err(err).Msg("Probe cannot be found:")
	} else if isAlreadyMirrored {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
