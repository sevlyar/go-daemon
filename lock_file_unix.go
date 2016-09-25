// +build darwin dragonfly freebsd linux netbsd openbsd plan9 solaris

package daemon

import (
	"fmt"
	"os"
	"syscall"
)

func lockFile(fd uintptr) error {
	err := syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK {
		err = ErrWouldBlock
	}
	return err
}

func unlockFile(fd uintptr) error {
	err := syscall.Flock(int(fd), syscall.LOCK_UN)
	if err == syscall.EWOULDBLOCK {
		err = ErrWouldBlock
	}
	return err
}

func getFdName(fd uintptr) (name string, err error) {
	// TODO(yar): This way does not work on darwin, use fcntl
	path := fmt.Sprintf("/proc/self/fd/%d", int(fd))

	var (
		fi os.FileInfo
		n  int
	)
	if fi, err = os.Lstat(path); err != nil {
		return
	}
	buf := make([]byte, fi.Size()+1)

	if n, err = syscall.Readlink(path, buf); err == nil {
		name = string(buf[:n])
	}
	return
}
