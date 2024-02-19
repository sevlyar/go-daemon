// +build !nacl

package daemon

import (
	"os"
	"syscall"
)

var handlers = make(map[os.Signal]SignalHandlerFunc)

func init() {
	handlers[syscall.SIGTERM] = sigtermDefaultHandler
}
