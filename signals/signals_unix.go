//go:build linux || bsd || darwin
// +build linux bsd darwin

package signals

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

func Wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	<-sigs
}
