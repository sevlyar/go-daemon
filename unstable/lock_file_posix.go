package daemon

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	// ErrWoldBlock indicates on locking pid-file by another process.
	ErrWouldBlock = syscall.EWOULDBLOCK
)

type LockFile struct {
	*os.File
}

func NewLockFile(file *os.File) *LockFile {
	return &LockFile{file}
}

func CreatePidFile(name string, perm os.FileMode) (lock *LockFile, err error) {
	if lock, err = OpenLockFile(name, perm); err != nil {
		return
	}
	if err = lock.Lock(); err != nil {
		lock.Remove()
		return
	}
	if err = lock.WritePid(); err != nil {
		lock.Remove()
	}
	return
}

func OpenLockFile(name string, perm os.FileMode) (lock *LockFile, err error) {
	var file *os.File
	if file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, perm); err == nil {
		lock = &LockFile{file}
	}
	return
}

func (file *LockFile) Lock() error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

func (file *LockFile) Unlock() error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

func (file *LockFile) WritePid() (err error) {
	if _, err = file.Seek(0, os.SEEK_SET); err != nil {
		return
	}
	var fileLen int
	if fileLen, err = fmt.Fprint(file, os.Getpid()); err != nil {
		return
	}
	if err = file.Truncate(int64(fileLen)); err != nil {
		return
	}
	err = file.Sync()
	return
}

func (file *LockFile) Remove() error {
	defer file.Close()

	if err := file.Unlock(); err != nil {
		log.Println(err)
		return err
	}

	name, err := GetFdName(file.Fd())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("daemons lock file name:", name)

	err = syscall.Unlink(name)
	return err
}

func GetFdName(fd uintptr) (name string, err error) {
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
