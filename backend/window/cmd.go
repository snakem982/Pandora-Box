package window

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

// RunCommandGetStdout executes a shell command and returns its stdout as a string.
// It works across Windows, macOS, and Linux.
func RunCommandGetStdout(command string, args ...string) (string, error) {
	var cmd *exec.Cmd

	// Special case: wrap in shell if needed (especially for Windows)
	if runtime.GOOS == "windows" {
		fullArgs := append([]string{"/C", command}, args...)
		cmd = exec.Command("cmd", fullArgs...)
	} else {
		cmd = exec.Command(command, args...)
	}

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stdoutBuf // Optional: capture stderr too

	err := cmd.Run()
	output := stdoutBuf.String()

	if err != nil {
		return output, fmt.Errorf("command failed: %w\nOutput: %s", err, output)
	}

	return output, nil
}
