package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var EnvVarName string = "GO_DAEMON"

func Daemonize(umask int, workDir string) (err error) {

	if isParent() {
		// parent process - fork and exec
		var path string
		if path, err = filepath.Abs(os.Args[0]); err != nil {
			return
		}

		cmd := prepareCommand(path)

		if err = cmd.Start(); err != nil {
			return
		}

		os.Exit(0)
	}

	// child process - daemon
	syscall.Umask(umask)

	if err = os.Chdir(workDir); err != nil {
		return
	}

	_, err = syscall.Setsid()

	// Do not requere redirect std
	// on /dev/null, this work done
	// forkExec func

	return
}

const DAEMON_VALUE = "1"

func isParent() bool {
	return os.Getenv(EnvVarName) != DAEMON_VALUE
}

func prepareCommand(path string) (cmd *exec.Cmd) {

	// prepare command-line arguments
	cmd = exec.Command(path, os.Args[1:]...)

	// prepare environment variables
	envVar := fmt.Sprintf("%s=%s", EnvVarName, DAEMON_VALUE)
	cmd.Env = append(os.Environ(), envVar)

	return
}
