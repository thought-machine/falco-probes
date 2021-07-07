package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/thought-machine/falco-probes/internal/cmd"
	"github.com/thought-machine/falco-probes/internal/logging"
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/resolver"
)

type opts struct {
	OutFile string `long:"out_file" description:"The path to a file to output the list of jobs for GitHub Actions" required:"yes"`
}

var log = logging.Logger

// Jobs represents a json-ifiable structure of jobs for producing a job matrix on GitHub Actions.
type Jobs []string

// JobsPerOperatingSystem returns Jobs per supported operating system.
func JobsPerOperatingSystem() Jobs {
	jobs := Jobs{}
	for os := range resolver.OperatingSystems {
		jobs = append(jobs, os)
	}

	return jobs
}

// WriteJobsToFile writes the given jobs to the given path as a JSON file.
func WriteJobsToFile(jobs Jobs, path string) error {
	jsonBytes, err := json.Marshal(jobs)
	if err != nil {
		return fmt.Errorf("could not marshal jobs: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	if err := ioutil.WriteFile(path, jsonBytes, 0644); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)

	jobs := JobsPerOperatingSystem()
	log.Info().
		Int("amount", len(jobs)).
		Msg("Generated jobs")

	if err := WriteJobsToFile(jobs, opts.OutFile); err != nil {
		log.Fatal().Err(err)
	}
	log.Info().
		Str("path", opts.OutFile).
		Msg("Wrote jobs to file")
}
