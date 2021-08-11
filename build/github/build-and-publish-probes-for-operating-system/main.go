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

// FalcoVersionNames represents the list of Falco versions to build eBPF probes for the given operating system. We're only interested in building the versions
// that diversify our support for Falco driver versions as they maintain compatibility between different Falco versions.
var FalcoVersionNames = []string{
	"0.24.0", // falco-driver-version: 85c88952b018fdbce2464222c3303229f5bfcfad
	"0.25.0", // falco-driver-version: ae104eb20ff0198a5dcb0c91cc36c86e7c3f25c7
	"0.26.0", // falco-driver-version: 2aa88dcf6243982697811df4c1b484bcbe9488a2
	"0.28.1", // falco-driver-version: 5c0b863ddade7a45568c0ac97d037422c9efb750
	"0.29.1", // falco-driver-version: 17f5df52a7d9ed6bb12d3b1768460def8439936d
}

type falcoVersion struct {
	Name   string
	Driver string
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

	log.Info().Msg("Getting list of falco drivers")
	FalcoVersions, err := getFalcoDrivers(cli, FalcoVersionNames)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get falco drivers")
	}
	log.Info().
		Int("amount", len(FalcoVersions)).
		Msg("Got falco drivers")

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
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		kernelPackageName := kernelPackageName

		parallelFns = append(parallelFns, func() error {
			return process1KernelPackage(
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

func process1KernelPackage(
	dockerCli *docker.Client,
	repo repository.Repository,
	operatingSystem operatingsystem.OperatingSystem,
	kernelPackageName string,
	falcoVersions []falcoVersion,
) error {
	// Get the required package specific values
	log.Info().
		Str("kernel_package_name", kernelPackageName).
		Msg("Getting kernel_package for")
	kernelPackage, err := operatingSystem.GetKernelPackageByName(kernelPackageName)
	if err != nil {
		return fmt.Errorf("could not get kernel package '%s': %w", kernelPackageName, err)
	}
	defer dockerCli.MustRemoveVolumes(
		kernelPackage.KernelSources,
		kernelPackage.KernelConfiguration,
	)
	probeName := kernelPackage.ProbeName() + ".o"
	log.Info().
		Str("kernel_package", kernelPackage.Name).
		Str("probe_name", probeName).
		Msg("Got kernel_package")

	for _, falcoVersion := range falcoVersions {
		// Check if probe is already mirrored to our repository & doesn't require building
		log.Info().
			Str("driver", falcoVersion.Driver).
			Str("probe_name", probeName).
			Msg("Checking whether probe is built & published")
		alreadyPublished, err := repo.IsAlreadyMirrored(falcoVersion.Driver, probeName)
		if err != nil {
			log.Error().Err(err).Msg("") // will just be logged as if probe is unfound it makes sense to try to build & publish it
		}

		if alreadyPublished {
			log.Info().
				Str("driver", falcoVersion.Driver).
				Str("probe_name", probeName).
				Msg("Skipping, probe is already built & published")
		} else {
			// Build unfound probe
			log.Info().
				Str("driver", falcoVersion.Driver).
				Str("probe_name", probeName).
				Msg("Not found, probe will now be built")
			builtDriverVersion, probePath, err := falcodriverbuilder.BuildEBPFProbe(
				dockerCli,
				falcoVersion.Name,
				operatingSystem,
				kernelPackage,
			)
			if err != nil {
				return fmt.Errorf("could not build probe for '%s': %w", kernelPackage.Name, err)
			}

			// Publish unfound probe
			log.Info().
				Str("driver", builtDriverVersion).
				Str("probe_path", probePath).
				Msg("Probe built, now publishing probe")
			if err := repo.PublishProbe(builtDriverVersion, probePath); err != nil {
				return fmt.Errorf("could not publish probe: %w", err)
			}
		}
	}

	return nil
}

func getFalcoDrivers(dockerCli *docker.Client, FalcoVersionNames []string) ([]falcoVersion, error) {
	var FalcoVersions []falcoVersion

	for _, falcoVersionName := range FalcoVersionNames {
		log.Info().
			Str("name", falcoVersionName).
			Msg("Getting driver for")
		falcoDriverBuilderImg, err := falcodriverbuilder.BuildImage(dockerCli, falcoVersionName)
		if err != nil {
			return FalcoVersions, fmt.Errorf("could not get falco_driver_builder_image for %s:%w", falcoVersionName, err)
		}
		driver, err := falcodriverbuilder.GetDriverVersion(dockerCli, falcoDriverBuilderImg)
		if err != nil {
			return FalcoVersions, fmt.Errorf("could not get driver for %s:%w", falcoDriverBuilderImg, err)
		}
		log.Info().
			Str("name", falcoVersionName).
			Str("driver", driver).
			Msg("Got driver")
		FalcoVersions = append(FalcoVersions, falcoVersion{falcoVersionName, driver})
	}

	return FalcoVersions, nil
}

func handleErrs(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Err(err).Msg("error encountered")
		}
		log.Fatal().Msg("errors encountered")
	}
}
