//go:build windows
// +build windows

package signals

import (
	"os"
	"os/signal"

	"golang.org/x/sys/windows"
)

func Wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, windows.SIGTERM, windows.SIGINT)
	<-sigs
}
