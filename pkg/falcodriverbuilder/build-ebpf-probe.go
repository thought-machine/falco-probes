package falcodriverbuilder

import (
	"fmt"

	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
)

var log = logging.Logger

// BuildEBPFProbe builds a Falco eBPF probe with the given falcoVersion, operatingsystem and kernelPackageName, returning the falcoDriverVersion and outProbePath.
func BuildEBPFProbe(
	cli *docker.Client,
	falcoVersion string,
	os operatingsystem.OperatingSystem,
	kernelPackage *operatingsystem.KernelPackage,
) (string, string, error) {
	log.Info().
		Str("falco_version", falcoVersion).
		Msg("Building falco-driver-builder")
	falcoDriverBuilderImage, err := BuildImage(cli, falcoVersion)
	if err != nil {
		return "", "", fmt.Errorf("could not build falco-driver-loader: %w", err)
	}
	falcoDriverVersion, err := GetDriverVersion(cli, falcoDriverBuilderImage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get falco driver version")
		return "", "", fmt.Errorf("could not get falco driver version: %w", err)
	}
	log.Info().
		Str("falco_version", falcoVersion).
		Str("falco_driver_version", falcoDriverVersion).
		Str("falco_driver_builder_image", falcoDriverBuilderImage).
		Msg("Built falco-driver-builder")

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
				builtProbeVolume:                  BuiltFalcoProbesDir,
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
		log.Error().Str("build-output", buildOut)
		return "", "", fmt.Errorf("could not build falco probe: %w", err)
	}

	builtProbePath, err := GetProbePathFromBuildOutput(buildOut)
	if err != nil {
		log.Error().Str("build-output", buildOut)
		return "", "", fmt.Errorf("could not build falco probe: %w", err)
	}

	probeReader, err := ExtractProbeFromVolume(cli, builtProbeVolume, builtProbePath)
	if err != nil {
		return "", "", fmt.Errorf("could not extract probe from built probe volume: %w", err)
	}

	outProbePath, err := WriteProbeToFile(falcoDriverVersion, builtProbePath, probeReader)
	if err != nil {
		return "", "", fmt.Errorf("could not write probe to file :%w", err)
	}

	log.Info().
		Str("path", outProbePath).
		Msg("successfully built probe")

	cli.MustRemoveVolumes(
		etcVolume,
		builtProbeVolume,
	)
	return falcoDriverVersion, outProbePath, nil
}
