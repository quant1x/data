//go:build !windows

package cross

import (
	"os"
	"syscall"
)

var StopSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP}
