package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

type opts struct {
	OutFile    string `long:"out_file" description:"The path to a file to output a list of Falco probes too (default: output to stdout)"`
	Positional struct {
		OperatingSystem string `positional-arg-name:"operating_system"`
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
		Msg("Getting kernel package names")
	kernelPackageNames, err := operatingSystem.GetKernelPackageNames()
	if err != nil {
		log.Fatal().Err(err).Msg("could not kernel package names")
	}

	log.Info().
		Int("amount", len(kernelPackageNames)).
		Str("operating_system", opts.Positional.OperatingSystem).
		Msg("got kernel packages")

	output := strings.Join(kernelPackageNames, "\n")

	if len(opts.OutFile) > 0 {
		ioutil.WriteFile(opts.OutFile, []byte(output), 0644)
		log.Info().
			Str("path", opts.OutFile).
			Msg("wrote kernel package names to file")

		return
	}

	fmt.Println(output)
}
