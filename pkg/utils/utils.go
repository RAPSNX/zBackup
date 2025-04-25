package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

func CMDWithOuput(name string, args ...string) (string, error) {
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

func CMDPipeOut(name string, args ...string) (io.ReadCloser, error) {
	cmd := exec.Command(name, args...)

	// Capture error
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	// Get stdout pipe
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"Error executing: %s %s, returnd: %w\n stderr: %s",
			name,
			strings.Join(args, " "),
			err,
			errBuf.String(),
		)
	}

	return stdoutPipe, nil
}

func CMDPipeIn(name string, input io.ReadCloser, args ...string) error {
	cmd := exec.Command(name, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(
			"Error executing: %s %s, returnd: %w\n stderr",
			name,
			strings.Join(args, " "),
			err,
		)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error reading stdout", err)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
