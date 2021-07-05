package main

import (
	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
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

	log.Info().
		Str("falco_version", opts.FalcoVersion).
		Msg("Building falco-driver-builder")
	falcoDriverBuilderImage, err := falcodriverbuilder.BuildImage(cli, opts.FalcoVersion)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build falco-driver-builder")
	}
	falcoDriverVersion, err := falcodriverbuilder.GetDriverVersion(cli, falcoDriverBuilderImage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get falco driver version")
	}
	log.Info().
		Str("falco_version", opts.FalcoVersion).
		Str("falco_driver_version", falcoDriverVersion).
		Str("falco_driver_builder_image", falcoDriverBuilderImage).
		Msg("Built falco-driver-builder")

	log.Info().
		Str("operating_system", opts.Positional.OperatingSystem).
		Str("kernel_package", opts.Positional.KernelPackage).
		Msg("Getting kernel package")
	kernelPackage, err := operatingSystem.GetKernelPackageByName(opts.Positional.KernelPackage)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("operating_system", opts.Positional.OperatingSystem).
			Str("kernel_package", opts.Positional.KernelPackage).
			Msg("could not get kernel package")
	}

	log.Info().Msg("Preparing /etc/os-release")
	etcVolume := cli.MustCreateVolume()
	if err := cli.WriteFileToVolume(etcVolume, "/etc/", "/etc/os-release", string(kernelPackage.OSRelease)); err != nil {
		log.Fatal().Err(err).Msg("could not write /etc/os-release")
	}

	log.Info().
		Str("operating_system", kernelPackage.OperatingSystem).
		Str("kernel_package", kernelPackage.Name).
		Str("kernel_release", kernelPackage.KernelRelease).
		Str("kernel_version", kernelPackage.KernelVersion).
		Str("kernel_machine", kernelPackage.KernelMachine).
		Str("falco_driver_version", falcoDriverVersion).
		Msg("Compiling Falco eBPF probe")
	builtProbeVolume := cli.MustCreateVolume()
	buildOut, err := cli.Run(
		&docker.RunOpts{
			Image: falcoDriverBuilderImage,
			Volumes: map[operatingsystem.Volume]string{
				builtProbeVolume:                  falcodriverbuilder.BuiltFalcoProbesDir,
				etcVolume:                         "/host/etc/",
				kernelPackage.KernelConfiguration: "/host/lib/modules/",
				kernelPackage.KernelSources:       "/host/usr/src/",
			},
			Env: map[string]string{
				"UNAME_V": kernelPackage.KernelVersion,
				"UNAME_R": kernelPackage.KernelRelease,
				"UNAME_M": kernelPackage.KernelMachine,
			},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Str("build-output", buildOut).Msg("could not build falco probe")
	}

	probeName, err := falcodriverbuilder.GetProbeNameFromBuildOutput(buildOut)
	if err != nil {
		log.Fatal().Err(err).Str("build-output", buildOut).Msg("could not find falco probe in build output")
	}

	probeReader, err := falcodriverbuilder.ExtractProbeFromVolume(cli, builtProbeVolume, probeName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not extract probe from built probe volume")
	}

	outProbePath, err := falcodriverbuilder.WriteProbeToFile(falcoDriverVersion, probeName, probeReader)
	if err != nil {
		log.Fatal().Err(err).Msg("could not write probe to file")
	}

	log.Info().
		Str("path", outProbePath).
		Msg("successfully built probe")

	cli.MustRemoveVolumes(
		kernelPackage.KernelSources,
		kernelPackage.KernelConfiguration,
		etcVolume,
		builtProbeVolume,
	)
}
