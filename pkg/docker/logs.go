package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	"github.com/docker/docker/pkg/jsonmessage"
)

// handleContainerLogs handles the container logs output returned by Docker and parses them into a string.
// This also additionally supports streaming the logs into a given extraWriter if set.
func handleContainerLogs(
	body io.ReadCloser,
	extraWriter io.Writer,
) string {
	scanner := bufio.NewScanner(body)

	out := bytes.Buffer{}
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) <= 8 {
			continue
		}
		// The first 8 bytes are a header as described in
		// https://github.com/moby/moby/issues/7375#issuecomment-51462963
		// so we can strip them out.
		line := b[8:]
		if bytes.ContainsRune(line, '\r') {
			// skip lines with carriage returns
			continue
		}

		if !bytes.HasSuffix(line, []byte("\n")) {
			// add a new line if there isn't one
			line = append(line, []byte("\n")...)
		}

		if extraWriter != nil {
			extraWriter.Write(line)
		}
		out.Write(line)
	}

	return out.String()
}

// handleBuildOrPullOutput handles the output of docker build or docker pull by parsing them and forwarding the parsed result into the given writer if set.
func handleBuildOrPullOutput(body io.ReadCloser, writer io.Writer) {
	decoder := json.NewDecoder(body)

	for decoder.More() {
		var bLog jsonmessage.JSONMessage
		if err := decoder.Decode(&bLog); err != nil {
			log.Warn().Err(err).Msg("could not unmarshal build output")
		}
		if writer != nil {
			writer.Write([]byte(bLog.Stream))
		}
	}
}
