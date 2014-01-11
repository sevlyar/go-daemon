package daemon

import (
	"log"
	"os"
	"syscall"
	"time"
)

func Example() {
	dmn := &Context{
		PidFileName: "/var/run/daemon.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "/",
		Umask:       027,
	}

	child, err := dmn.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if child != nil {
		return
	}
	defer dmn.Release()

	go func() {
		for {
			time.Sleep(0)
		}
	}()

	termHandler := func(sig os.Signal) error {
		log.Println("SIGTERM:", sig)
		return ErrStop
	}

	hupHandler := func(sig os.Signal) error {
		log.Println("SIGHUP:", sig)
		return nil
	}

	SetSigHandler(termHandler, syscall.SIGTERM, syscall.SIGKILL)
	SetSigHandler(hupHandler, syscall.SIGHUP)

	err = ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
}
