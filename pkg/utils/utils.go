package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func SanitizeRawStringToList[T ~string](raw string) []T {
	var sanitized []T
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") {
			sanitized = append(sanitized, T(line))
		}
	}
	return sanitized
}

func RunCMD(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	// Capture output
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"Error executing: %s %s, returnd: %w\n stderr: %s",
			name,
			strings.Join(args, " "),
			err,
			errBuf.String(),
		)
	}

	return outBuf.String(), nil
}
