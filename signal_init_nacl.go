// +build nacl

package daemon

import (
	"os"
)

var handlers = make(map[os.Signal]SignalHandlerFunc)

// func init() {
// 	handlers[syscall.SIGTERM] = sigtermDefaultHandler
// }
