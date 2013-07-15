package main

import (
	"errors"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	daemon.Reborn(027, "./")

	file, _ := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	daemon.RedirectStream(os.Stdout, file)
	daemon.RedirectStream(os.Stderr, file)
	file.Close()
	log.Println("--- log ---")

	SignalsHandler(TermHandler, syscall.SIGTERM, syscall.SIGKILL)
	SignalsHandler(HupHandler, syscall.SIGHUP)
	SignalsHandler(Usr1Handler, syscall.SIGUSR1)
	if err := ServeSignals(); err != nil {
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

type SignalHandlerFunc func(sig os.Signal) (stop bool, err error)

func SignalsHandler(handler SignalHandlerFunc, signals ...os.Signal) {
	for _, sig := range signals {
		handlers[sig] = handler
	}
}

func ServeSignals() (err error) {
	signals := make([]os.Signal, 0, len(handlers))
	for sig, _ := range handlers {
		signals = append(signals, sig)
	}

	ch := make(chan os.Signal, 8)
	signal.Notify(ch, signals...)

	var stop bool
	for sig := range ch {
		stop, err = handlers[sig](sig)
		if stop {
			break
		}
	}

	signal.Stop(ch)

	return
}

var handlers = make(map[os.Signal]SignalHandlerFunc)
