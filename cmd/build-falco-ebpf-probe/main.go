package main

import (
	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

type opts struct {
	FalcoVersion string `long:"falco_version" description:"The version of Falco to compile probes against" required:"true"`
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

	log.Info().
		Str("operating_system", opts.Positional.OperatingSystem).
		Msg("Resolving operating system")
	operatingSystem, err := resolver.OperatingSystem(cli, opts.Positional.OperatingSystem)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get operating system")
	}

	kernelPackage, err := operatingSystem.GetKernelPackageByName(opts.Positional.KernelPackage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get kernel package")
	}

	if _, _, err := falcodriverbuilder.BuildEBPFProbe(
		cli,
		opts.FalcoVersion,
		operatingSystem,
		kernelPackage,
	); err != nil {
		log.Fatal().Err(err).Msg("could not build eBPF probe")
	}

	cli.MustRemoveVolumes(
		kernelPackage.KernelSources,
		kernelPackage.KernelConfiguration,
	)
}
