package daemon

import (
	"errors"
	"os"
	"os/signal"
)

var ErrStop = errors.New("stop serve signals")

type SignalHandlerFunc func(sig os.Signal) (err error)

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

	if err == ErrStop {
		err = nil
	}

	return
}

var handlers = make(map[os.Signal]SignalHandlerFunc)
