//go:build !windows

package isadmin

import "os"

// Check if the program has administrative privileges.
func Check() bool {
	return os.Getuid() == 0
}
