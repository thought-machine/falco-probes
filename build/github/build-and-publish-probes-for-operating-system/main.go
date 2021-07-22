package main

import (
	"fmt"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
	"github.com/thought-machine/falco-probes/pkg/repository"
	"github.com/thought-machine/falco-probes/pkg/repository/ghreleases"
)

var log = logging.Logger

// FalcoVersions represents the list of Falco versions to build eBPF probes for the given operating system. We're only interested in building the versions
// that diversify our support for Falco driver versions as they maintain compatibility between different Falco versions.
// Note: We can only support 0.28.1+ at the moment as it seems like the falco-driver-loader script changed in an incompatible way between 0.26 and 0.28.1.
// TODO: To fix this, we could just source the driver loader script from 0.28.1 and reuse that instead of the script bundled w/ each falco-driver-loader.
var FalcoVersions = []string{
	// "0.24.0", // falco-driver-version: 85c88952b018fdbce2464222c3303229f5bfcfad
	// "0.25.0", // falco-driver-version: ae104eb20ff0198a5dcb0c91cc36c86e7c3f25c7
	// "0.26.0", // falco-driver-version: 2aa88dcf6243982697811df4c1b484bcbe9488a2
	"0.28.1", // falco-driver-version: 5c0b863ddade7a45568c0ac97d037422c9efb750
	"0.29.1", // falco-driver-version: 17f5df52a7d9ed6bb12d3b1768460def8439936d
}

type opts struct {
	Parallelism int             `long:"parallelism" description:"The amount of probes to compile at the same time" default:"4"`
	GHReleases  ghreleases.Opts `group:"github_releases" namespace:"github_releases"`
	Positional  struct {
		OperatingSystem string `positional-arg-name:"operating_system"`
	} `positional-args:"yes" required:"true"`
}

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)

	cli := docker.MustClient()
	ghReleases := ghreleases.MustGHReleases(&opts.GHReleases)

	log.Info().
		Str("operating_system", opts.Positional.OperatingSystem).
		Msg("Resolving operating system")
	operatingSystem, err := resolver.OperatingSystem(cli, opts.Positional.OperatingSystem)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get operating system")
	}

	log.Info().
		Str("operating_system", opts.Positional.OperatingSystem).
		Msg("Getting list of kernel packages")
	kernelPackageNames, err := operatingSystem.GetKernelPackageNames()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get kernel package names")
	}

	log.Info().
		Int("amount", len(kernelPackageNames)).
		Msg("Retrieving kernel packages")

	parallelFns := []func() error{}
	for _, kernelPackageName := range kernelPackageNames {
		parallelFns = append(parallelFns, func() error {
			return buildProbesForKernelPackageName(
				cli,
				ghReleases,
				operatingSystem,
				kernelPackageName,
				FalcoVersions,
			)
		})
	}

	errs := cmd.RunParallelAndCollectErrors(parallelFns, opts.Parallelism)

	handleErrs(errs)
}

func buildProbesForKernelPackageName(
	dockerCli *docker.Client,
	repo repository.Repository,
	operatingSystem operatingsystem.OperatingSystem,
	kernelPackageName string,
	falcoVersions []string,
) error {
	log.Info().
		Str("kernel_package", kernelPackageName).
		Msg("Getting kernel package")

	kernelPackage, err := operatingSystem.GetKernelPackageByName(kernelPackageName)
	if err != nil {
		return fmt.Errorf("could not get kernel package '%s': %w", kernelPackageName, err)
	}
	defer dockerCli.MustRemoveVolumes(
		kernelPackage.KernelSources,
		kernelPackage.KernelConfiguration,
	)

	log.Info().
		Str("kernel_package", kernelPackage.Name).
		Msg("Got kernel package")

	for _, falcoVersion := range falcoVersions {
		log.Info().
			Str("kernel_package", kernelPackage.Name).
			Str("falco_version", falcoVersion).
			Msg("Building Falco eBPF probe")

		falcoDriverVersion, probePath, err := falcodriverbuilder.BuildEBPFProbe(
			dockerCli,
			falcoVersion,
			operatingSystem,
			kernelPackage,
		)
		if err != nil {
			return fmt.Errorf("could not build eBPF probe for '%s': %w", kernelPackage.Name, err)
		}

		if err := repo.PublishProbe(falcoDriverVersion, probePath); err != nil {
			return fmt.Errorf("could not publish probe: %w", err)
		}
	}

	return nil
}

func handleErrs(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Err(err).Msg("error encountered")
		}
		log.Fatal().Msg("errors encountered")
	}
}
