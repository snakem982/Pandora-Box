package utils

import (
	"fmt"
	"github.com/gofrs/flock"
)

var fileLock *flock.Flock

// NotSingleton attempts to acquire a lock file to ensure only one instance is running.
// Returns true if another instance is already running.
func NotSingleton(name string) bool {
	lockFile := GetUserHomeDir("pid", name)
	file, err := CreateFile(lockFile)
	if err == nil && file != nil {
		_ = file.Close()
	}

	fileLock = flock.New(lockFile)

	locked, err := fileLock.TryLock()
	if err != nil {
		fmt.Printf("Failed to acquire lock for %s: %v\n", name, err)
		return true
	}

	if !locked {
		fmt.Printf("Another instance of %s is already running.\n", name)
		return true
	}

	return false
}

// UnlockSingleton should be called before exit to release the lock (optional).
func UnlockSingleton() {
	if fileLock != nil {
		_ = fileLock.Unlock()
	}
}
