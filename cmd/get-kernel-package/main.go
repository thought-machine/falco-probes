package main

import (
	"os"
	"path/filepath"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

type opts struct {
	OutputDir  string `long:"output_dir" description:""`
	Positional struct {
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
		Str("operating_system", opts.Positional.OperatingSystem).
		Str("kernel_package_name", opts.Positional.KernelPackage).
		Msg("Getting kernel package")
	kernelPackage, err := operatingSystem.GetKernelPackageByName(opts.Positional.KernelPackage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get kernel package")
	}

	if err := os.MkdirAll(opts.OutputDir, 0777); err != nil {
		log.Fatal().Err(err).Msg("could not create output dir")
	}

	if err := os.WriteFile(filepath.Join(opts.OutputDir, "kernel_release"), []byte(kernelPackage.KernelRelease), 0644); err != nil {
		log.Fatal().Err(err).Msg("could not write file")
	}

	if err := os.WriteFile(filepath.Join(opts.OutputDir, "kernel_version"), []byte(kernelPackage.KernelVersion), 0644); err != nil {
		log.Fatal().Err(err).Msg("could not write file")
	}

	if err := os.WriteFile(filepath.Join(opts.OutputDir, "kernel_machine"), []byte(kernelPackage.KernelMachine), 0644); err != nil {
		log.Fatal().Err(err).Msg("could not write file")
	}

	if err := os.WriteFile(filepath.Join(opts.OutputDir, "os_release"), []byte(kernelPackage.OSRelease), 0644); err != nil {
		log.Fatal().Err(err).Msg("could not write file")
	}

	if err := cli.GetDirectoryFromVolume(
		kernelPackage.KernelConfiguration,
		"/lib/modules/",
		"/lib/modules/",
		filepath.Join(opts.OutputDir, "kernel_configuration"),
	); err != nil {
		log.Fatal().Err(err).Msg("could not get kernel configuration")
	}

	if err := cli.GetDirectoryFromVolume(
		kernelPackage.KernelSources,
		"/usr/src/",
		"/usr/src/",
		filepath.Join(opts.OutputDir, "kernel_sources"),
	); err != nil {
		log.Fatal().Err(err).Msg("could not get kernel sources")
	}

}
