package daemon

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestCompilation(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	if !requireMinor(5) {
		t.Skip(runtime.Version(), "cross-compilation requires compiler bootstrapping")
	}

	pairs := []string{
		"darwin/386",
		"darwin/amd64",
		"dragonfly/amd64",
		"freebsd/386",
		"freebsd/amd64",
		"freebsd/arm",
		"linux/386",
		"linux/amd64",
		"linux/arm",
		"linux/arm64",
		"netbsd/386",
		"netbsd/amd64",
		"netbsd/arm",
		"openbsd/386",
		"openbsd/amd64",
		"openbsd/arm",
		"solaris/amd64",
		"windows/386",
		"windows/amd64",

		// TODO(yar): support plan9
		//"plan9/386",
		//"plan9/amd64",
		//"plan9/arm",
	}

	env := os.Environ()
	for i := range pairs {
		p := pairs[i]
		pair := strings.Split(p, "/")
		goos, goarch := pair[0], pair[1]
		if goos == "solaris" && !requireMinor(7) {
			t.Log("skip, solaris requires at least go1.7")
			continue
		}
		cmd := exec.Command("go", "build", "./")
		env := append([]string(nil), env...)
		cmd.Env = append(env, "GOOS="+goos, "GOARCH="+goarch)
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			t.Log(p, "\n", string(out))
		}
		if err != nil {
			t.Error(p, err)
		}
	}
}

func requireMinor(minor int) bool {
	str := runtime.Version()
	if !strings.HasPrefix(str, "go1.") {
		return true
	}
	str = strings.TrimPrefix(str, "go1.")
	ver, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return ver >= minor
}
