package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	envVarName  = "_GO_DAEMON"
	envVarValue = "1"
)

// func Reborn daemonize process.
func Reborn(umask uint32, workDir string) (err error) {

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
	syscall.Umask(int(umask))

	if len(workDir) == 0 {
		if err = os.Chdir(workDir); err != nil {
			return
		}
	}

	_, err = syscall.Setsid()

	// Do not requere redirect std
	// on /dev/null, this work done
	// forkExec func

	return
}

func isParent() bool {
	return os.Getenv(envVarName) != envVarValue
}

func IsWasReborn() bool {
	return !isParent()
}

func prepareCommand(path string) (cmd *exec.Cmd) {

	// prepare command-line arguments
	cmd = exec.Command(path, os.Args[1:]...)

	// prepare environment variables
	envVar := fmt.Sprintf("%s=%s", envVarName, envVarValue)
	cmd.Env = append(os.Environ(), envVar)

	return
}

func RedirectStream(stream, target *os.File) (err error) {

	stdoutFd := int(stream.Fd())
	if err = syscall.Close(stdoutFd); err != nil {
		return
	}

	err = syscall.Dup2(int(target.Fd()), stdoutFd)

	return
}
