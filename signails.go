package daemon

import (
	"os"
	"os/signal"
)

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
