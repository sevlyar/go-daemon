package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

const (
	markName  = "_GO_DAEMON"
	markValue = "1"
)

func WasReborn() bool {
	return os.Getenv(markName) == markValue
}

type Context struct {
	PidFile string
	PidPerm os.FileMode
	LogFile string
	LogPerm os.FileMode

	WorkDir    string
	Chroot     string
	Env        []string
	Args       []string
	Credential *syscall.Credential
	Umask      int

	abspath  string
	pidFile  *os.File
	logFile  *os.File
	nullFile *os.File

	rpipe, wpipe *os.File

	hidden *Context
}

func (d *Context) Reborn() (child *os.Process, err error) {
	if !WasReborn() {
		child, err = d.parent()
	} else {
		err = d.child()
	}
	return
}

func (d *Context) openFiles() (err error) {
	if d.nullFile, err = os.Open(os.DevNull); err != nil {
		return
	}

	if len(d.PidFile) > 0 {
		var pid *LockFile
		if pid, err = CreatePidFile(d.PidFile, d.PidPerm); err != nil {
			return
		}
		d.pidFile = pid.File
	} else {
		d.pidFile = d.nullFile
	}

	if len(d.LogFile) > 0 {
		if d.logFile, err = os.OpenFile(d.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, d.LogPerm); err != nil {
			return
		}
	} else {
		d.logFile = d.nullFile
	}

	d.rpipe, d.wpipe, err = os.Pipe()
	return
}

func (d *Context) closeFiles() (err error) {
	return
}

func (d *Context) prepareEnv() (err error) {
	if d.abspath, err = filepath.Abs(os.Args[0]); err != nil {
		return
	}

	if len(d.Args) == 0 {
		d.Args = os.Args[1:]
	}

	mark := fmt.Sprintf("%s=%s", markName, markValue)
	if len(d.Env) == 0 {
		d.Env = os.Environ()
	}
	d.Env = append(d.Env, mark)

	return
}

func (d *Context) parent() (child *os.Process, err error) {
	if err = d.prepareEnv(); err != nil {
		return
	}

	if err = d.openFiles(); err != nil {
		return
	}
	defer d.closeFiles()

	attr := &os.ProcAttr{
		Dir:   d.WorkDir,
		Env:   d.Env,
		Files: []*os.File{d.rpipe, d.nullFile, d.logFile, d.pidFile},
		Sys: &syscall.SysProcAttr{
			Chroot:     d.Chroot,
			Credential: d.Credential,
			Setsid:     true,
		},
	}

	if child, err = os.StartProcess(d.abspath, d.Args, attr); err != nil {
		return
	}

	d.rpipe.Close()
	encoder := json.NewEncoder(d.wpipe)
	err = encoder.Encode(d)

	return
}

func (d *Context) child() (err error) {
	//d.hidden = new(Context)
	decoder := json.NewDecoder(os.Stdin)
	if err = decoder.Decode(d); err != nil {
		return
	}

	syscall.Umask(int(d.Umask))

	// TODO: replace /dev/null
	DupFile(os.Stdin, os.Stdout)

	return
}

func (d *Context) Release() {

}

// TODO: rename this
// func RedirectStream redirects file s to file target.
func DupFile(s, target *os.File) (err error) {

	stdoutFd := int(s.Fd())
	if err = syscall.Close(stdoutFd); err != nil {
		return
	}

	err = syscall.Dup2(int(target.Fd()), stdoutFd)

	return
}
