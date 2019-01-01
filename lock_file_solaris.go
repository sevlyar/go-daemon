// +build solaris

package daemon

import (
	"golang.org/x/sys/unix"
	"os"
)

func lockFile(fd uintptr) error {

	flockT := unix.Flock_t{
		Type:   unix.F_WRLCK,
		Whence: int16(os.SEEK_SET),
		Start:  0,
		Len:    0,
		Pid:    0,
	}

	err := unix.FcntlFlock(fd, unix.F_SETLK, &flockT)
	return err
}

func unlockFile(fd uintptr) error {

	flockT := unix.Flock_t{
		Type:   unix.F_WRLCK,
		Whence: int16(os.SEEK_SET),
		Start:  0,
		Len:    0,
		Pid:    0,
	}

	err := unix.FcntlFlock(fd, unix.F_SETLKW, &flockT)
	return err
}
