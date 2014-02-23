package daemon

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var (
	filename                = os.TempDir() + "/test.lock"
	fileperm    os.FileMode = 0644
	invalidname             = "/x/y/unknown"
)

func TestCreatePidFile(test *testing.T) {
	if _, err := CreatePidFile(invalidname, fileperm); err == nil {
		test.Fatal("CreatePidFile(): Error was not detected on invalid name")
	}

	lock, err := CreatePidFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	defer lock.Remove()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		test.Fatal(err)
	}
	if string(data) != fmt.Sprint(os.Getpid()) {
		test.Fatal("pids not equal")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	if err = NewLockFile(file).WritePid(); err == nil {
		test.Fatal("WritePid(): Error was not detected on invalid permissions")
	}
}

func TestNewLockFile(test *testing.T) {
	lock := NewLockFile(os.NewFile(1001, ""))
	err := lock.Remove()
	if err == nil {
		test.Fatal("Remove(): Error was not detected on invalid fd")
	}
	err = lock.WritePid()
	if err == nil {
		test.Fatal("WritePid(): Error was not detected on invalid fd")
	}
}

func TestGetFdName(test *testing.T) {
	name, err := GetFdName(0)
	if err != nil {
		test.Error(err)
	} else {
		if name != "/dev/null" {
			test.Errorf("Filename of fd 0: `%s'", name)
		}
	}

	name, err = GetFdName(1011)
	if err == nil {
		test.Errorf("GetFdName(): Error was not detected on invalid fd, name: `%s'", name)
	}
}

func TestReadPid(test *testing.T) {
	lock, err := CreatePidFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	defer lock.Remove()

	pid, err := lock.ReadPid()
	if err != nil {
		test.Fatal("ReadPid(): Unable read pid from file:", err)
	}

	if pid != os.Getpid() {
		test.Fatal("Pid not equal real pid")
	}
}

func TestLockFileLock(test *testing.T) {
	lock, err := OpenLockFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	defer lock.Remove()

	if err := lock.Lock(); err != nil {
		test.Fatal(err)
	}
	scr, msg, err := createLockScriptAndStart()
	if err != nil {
		test.Fatal(err)
	}
	if msg != "error" {
		test.Fatal("script was able lock file")
	}
	if err = terminateLockScript(scr); err != nil {
		test.Fatal(err)
	}

	if err = lock.Unlock(); err != nil {
		test.Fatal(err)
	}
	lock.Close()

	scr, msg, err = createLockScriptAndStart()
	if err != nil {
		test.Fatal(err)
	}
	if msg != "locked" {
		test.Fatal("script can not lock file")
	}
	lock, err = CreatePidFile(filename, fileperm)
	if err != ErrWouldBlock {
		test.Fatal("Lock() not work properly")
	}
	if err = terminateLockScript(scr); err != nil {
		test.Error(err)
	}
}

func createLockScriptAndStart() (scr *script, msg string, err error) {
	var text = fmt.Sprintf(`
		set -e
		exec 222>"%s"
		flock -n 222||echo "error"
		echo "locked"
		read inp`, filename)

	scr, err = createScript(text, true)
	if err != nil {
		return
	}

	if err = scr.cmd.Start(); err != nil {
		return
	}
	// wait until the script does not try to lock the file
	msg, err = scr.get()
	return
}

func terminateLockScript(scr *script) (err error) {
	if err = scr.send(""); err != nil {
		return
	}
	err = scr.cmd.Wait()
	return
}

type script struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stdin  io.WriteCloser
}

func createScript(text string, createPipes bool) (scr *script, err error) {
	var scrName string
	if scrName, err = createScriptFile(text); err != nil {
		return
	}
	scr = &script{cmd: exec.Command("bash", scrName)}
	if createPipes {
		if scr.stdout, err = scr.cmd.StdoutPipe(); err != nil {
			return
		}
		if scr.stdin, err = scr.cmd.StdinPipe(); err != nil {
			return
		}
	}
	return
}

func (scr *script) send(line string) (err error) {
	_, err = fmt.Fprintln(scr.stdin, line)
	return
}

func (scr *script) get() (line string, err error) {
	_, err = fmt.Fscanln(scr.stdout, &line)
	return
}

func createScriptFile(text string) (name string, err error) {
	var scr *os.File
	if scr, err = ioutil.TempFile(os.TempDir(), "scr"); err != nil {
		return
	}
	defer scr.Close()
	name = scr.Name()
	_, err = scr.WriteString(text)
	return
}
