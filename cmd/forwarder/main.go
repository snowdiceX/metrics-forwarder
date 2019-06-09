package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/snowdiceX/metrics-forwarder/log"
)

func main() {
	defer log.Flush()

	root := NewRootCommand()
	root.AddCommand(
		NewStartCommand(starter, true),
		NewVersionCommand(versioner, false))

	if err := root.Execute(); err != nil {
		log.Error("Exit by error: ", err)
	}
}

// KeepRunning just keep running
func KeepRunning(callback func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case s, ok := <-signals:
		log.Infof("System signal [%v] %t, trying to run callback...", s, ok)
		if !ok {
			break
		}
		if callback != nil {
			callback(s)
		}
		log.Flush()
		os.Exit(1)
	}
}
