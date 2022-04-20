package util

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

var shutdownRequests = make(chan bool)

type ShutdownTimeouts struct {
	KillTimeout      time.Duration
	InterruptTimeout time.Duration
}

func Shutdown() {
	shutdownRequests <- true
}

func GetShutdownRequests() <-chan bool {
	return shutdownRequests
}

func SetupShutdownOnSignals(timeouts ShutdownTimeouts, signals ...os.Signal) <-chan bool {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)

	go func() {
		for { // Drain the channel & update the shutdownRequests channel
			sig := <-sigs

			Shutdown()
			switch sig {
			case os.Interrupt:
				go exitSignalWithTimeout("Interrupt", timeouts.InterruptTimeout)
			case os.Kill:
				go exitSignalWithTimeout("Kill", timeouts.KillTimeout)
			}
		}
	}()

	return shutdownRequests
}

func exitSignalWithTimeout(signal string, timeout time.Duration) {
	log.Error().Msgf("%s signal received.", signal)
	if timeout > 0 {
		time.Sleep(timeout)
		log.Error().Msgf("%s signal timeout (%v) expired. Exiting...", signal, timeout)
		os.Exit(1) // TODO: Check if this is the correct exit code
	}
}

func SetupShutdown(timeouts ...ShutdownTimeouts) <-chan bool {
	var timeout ShutdownTimeouts
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	} else {
		timeout = ShutdownTimeouts{}
	}

	return SetupShutdownOnSignals(timeout, os.Interrupt, os.Kill)
}
