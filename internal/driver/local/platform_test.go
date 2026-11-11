package local

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// sleeper is a process that only waits, for about the given seconds: sleep, or ping on Windows,
// which has no sleep program.
func sleeper(seconds int) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("ping", "-n", strconv.Itoa(seconds), "127.0.0.1")
	}
	return exec.Command("sleep", strconv.Itoa(seconds))
}

// needSymlinks skips a test whose fixture is a symlink where this user may not create one
// (Windows without Developer Mode).
func needSymlinks(t testing.TB) {
	t.Helper()
	if err := os.Symlink("target", filepath.Join(t.TempDir(), "link")); err != nil {
		t.Skip("this user may not create symlinks: " + err.Error())
	}
}

// skipOnWindows skips a test that needs a POSIX shell or signals.
func skipOnWindows(t testing.TB, why string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("not on Windows: " + why)
	}
}
