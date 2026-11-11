package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// needSymlinks skips a test whose fixture is a symlink where this user may not create one
// (Windows without Developer Mode).
func needSymlinks(t testing.TB) {
	t.Helper()
	if err := os.Symlink("target", filepath.Join(t.TempDir(), "link")); err != nil {
		t.Skip("this user may not create symlinks: " + err.Error())
	}
}

// fakeProgram puts this test binary in dir under name, for PATH to find: a symlink, or a copy
// where this user may not create one (a hard link to the running binary could not be removed
// with the test's directory).
func fakeProgram(t testing.TB, dir, name string) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, name+exeSuffix)
	if os.Symlink(exe, link) == nil {
		return
	}
	b, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(link, b, 0o700); err != nil {
		t.Fatal(err)
	}
}

// exeSuffix is what a program's file name needs for PATH to find it.
var exeSuffix = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

// shortTemp is where a test's piggery home goes: /tmp keeps unix socket paths short (macOS limits
// them); Windows has no /tmp, and its daemon pipe name does not depend on the path's length.
func shortTemp() string {
	if runtime.GOOS == "windows" {
		return ""
	}
	return "/tmp"
}

// skipOnWindows skips a test that needs a POSIX shell or signals.
func skipOnWindows(t testing.TB, why string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("not on Windows: " + why)
	}
}

// unixModes: file permission bits mean something (not on Windows, where a file is 0666 or 0444).
var unixModes = runtime.GOOS != "windows"
