package main

import (
	"context"
	"sync"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/internal/queue"
	"github.com/thought-machine/falco-probes/pkg/docker"
	"github.com/thought-machine/falco-probes/pkg/falcodriverbuilder"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

// FalcoVersionNames represents the list of Falco versions to build eBPF probes for the given operating system. We're only interested in building the versions
// that diversify our support for Falco driver versions as they maintain compatibility between different Falco versions.
var FalcoVersionNames = []string{
	"0.24.0", // falco-driver-version: 85c88952b018fdbce2464222c3303229f5bfcfad
	"0.25.0", // falco-driver-version: ae104eb20ff0198a5dcb0c91cc36c86e7c3f25c7
	"0.26.0", // falco-driver-version: 2aa88dcf6243982697811df4c1b484bcbe9488a2
	"0.28.1", // falco-driver-version: 5c0b863ddade7a45568c0ac97d037422c9efb750
	"0.29.1", // falco-driver-version: 17f5df52a7d9ed6bb12d3b1768460def8439936d
}

var log = logging.Logger

type opts struct {
	Parallelism int    `long:"parallelism" description:"The amount of probes to compile at the same time" default:"4"`
	Buffer      uint64 `long:"buffer" description:"The maximum amount of tasks to buffer." default:"128"`
	Positional  struct {
		OperatingSystem string `positional-arg-name:"operating_system"`
	} `positional-args:"yes" required:"true"`
}

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)

	docker.DockerClient.MustCleanupVolumes()

	operatingSystem, err := resolver.OperatingSystem(docker.DockerClient, opts.Positional.OperatingSystem)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get operating system")
	}

	q := queue.NewQueue(&queue.Opts{
		Buffer: opts.Buffer,
	})

	ctx := context.Background()
	ctx, cancelFn := context.WithCancel(ctx)
	_ = cancelFn

	var wg sync.WaitGroup
	for i := 0; i < opts.Parallelism; i++ {
		wg.Add(1)
		go queue.Worker(ctx, &wg, q)
	}

	q.Publish(&GetKernelPackageNamesTask{
		OperatingSystem: operatingSystem,
	})

	wg.Wait()
}

type GetKernelPackageNamesTask struct {
	queue.Task
	OperatingSystem operatingsystem.OperatingSystem
}

func (t *GetKernelPackageNamesTask) Execute(queue *queue.Queue) error {
	defer queue.Ack()

	log.Info().Str("operating-system", t.OperatingSystem.GetName()).Msg("getting kernel package names")
	packageNames, err := t.OperatingSystem.GetKernelPackageNames()
	if err != nil {
		return err
	}

	log.Info().Str("operating-system", t.OperatingSystem.GetName()).Int("amount", len(packageNames)).Msg("got kernel package names")

	for _, packageName := range packageNames {
		log.Debug().
			Str("operating-system", t.OperatingSystem.GetName()).
			Str("name", packageName).
			Msg("adding task to get kernel package")

		queue.Publish(&GetKernelPackageByNameTask{
			OperatingSystem: t.OperatingSystem,
			Name:            packageName,
		})
	}

	return nil
}

type GetKernelPackageByNameTask struct {
	queue.Task
	OperatingSystem operatingsystem.OperatingSystem
	Name            string
}

func (t *GetKernelPackageByNameTask) Execute(queue *queue.Queue) error {
	defer queue.Ack()

	log.Info().Str("operating-system", t.OperatingSystem.GetName()).Str("name", t.Name).Msg("getting kernel package")

	kernelPackage, err := t.OperatingSystem.GetKernelPackageByName(t.Name)
	if err != nil {
		return err
	}

	log.Info().Str("operating-system", t.OperatingSystem.GetName()).Str("name", t.Name).Msg("got kernel package")

	for _, falcoVersion := range FalcoVersionNames {
		queue.Publish(&BuildEBPFProbeTask{
			OperatingSystem: t.OperatingSystem,
			FalcoVersion:    falcoVersion,
			KernelPackage:   kernelPackage,
		})
	}

	return nil
}

type BuildEBPFProbeTask struct {
	queue.Task
	OperatingSystem operatingsystem.OperatingSystem
	FalcoVersion    string
	KernelPackage   *operatingsystem.KernelPackage
}

func (t *BuildEBPFProbeTask) Execute(queue *queue.Queue) error {
	defer queue.Ack()

	log.Info().
		Str("falco-version", t.FalcoVersion).
		Str("operating-system", t.OperatingSystem.GetName()).
		Str("name", t.KernelPackage.Name).
		Msg("building eBPF probe")

	if _, _, err := falcodriverbuilder.BuildEBPFProbe(
		docker.DockerClient,
		t.FalcoVersion,
		t.OperatingSystem,
		t.KernelPackage,
	); err != nil {
		return err
	}

	return nil
}
