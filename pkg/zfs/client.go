package zfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type client struct{}

func (c *client) list() (string, error) {
	return execWithOut("zfs", "list")
}

func execWithOut(name string, args ...string) (string, error) {
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
