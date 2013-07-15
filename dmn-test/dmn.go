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
	ServeSignals()

	log.Println("--- end ---")
}

func TermHandler(sig os.Signal) error {
	log.Println("SIGTERM:", sig)
	return ErrServTerm
}

func HupHandler(sig os.Signal) error {
	log.Println("SIGHUP:", sig)
	return nil
}

type SignalHandlerFunc func(sig os.Signal) error

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

	for sig := range ch {
		err = handlers[sig](sig)
		if err != nil {
			break
		}
	}

	signal.Stop(ch)

	return
}

var ErrServTerm = errors.New("Termination signals service")

var handlers = make(map[os.Signal]SignalHandlerFunc)
