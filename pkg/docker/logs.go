package docker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// HandleContainerLogs handles the container logs output returned by Docker and parses them into a string.
func HandleContainerLogs(body io.ReadCloser) string {
	scanner := bufio.NewScanner(body)

	out := ""
	for scanner.Scan() {
		b := scanner.Bytes()
		out += handleContainerLog(b)
	}

	return out
}

func handleContainerLog(b []byte) string {
	if len(b) <= 8 {
		return ""
	}
	// The first 8 bytes are a header as described in
	// https://github.com/moby/moby/issues/7375#issuecomment-51462963
	// so we can strip them out.
	line := string(b[8:])

	if !strings.Contains(line, "\r") {
		return fmt.Sprintf("%s\n", line)
	}

	return line
}
