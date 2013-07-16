package main

import (
	"errors"
	"go-daemon"
	"log"
	"os"
	"syscall"
)

func main() {
	daemon.Reborn(027, "./")

	file, _ := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	daemon.RedirectStream(os.Stdout, file)
	daemon.RedirectStream(os.Stderr, file)
	file.Close()
	log.Println("--- log ---")

	daemon.SignalsHandler(TermHandler, syscall.SIGTERM, syscall.SIGKILL)
	daemon.SignalsHandler(HupHandler, syscall.SIGHUP)
	daemon.SignalsHandler(Usr1Handler, syscall.SIGUSR1)

	err := daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}

	log.Println("--- end ---")
}

func TermHandler(sig os.Signal) (stop bool, err error) {
	log.Println("SIGTERM:", sig)
	stop = true
	return
}

func HupHandler(sig os.Signal) (stop bool, err error) {
	log.Println("SIGHUP:", sig)
	stop = false
	return
}

func Usr1Handler(sig os.Signal) (stop bool, err error) {
	log.Println("SIGUSR1:", sig)
	return true, errors.New("some error")
}
