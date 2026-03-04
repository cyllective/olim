package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/cyllective/olim/internal/app"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_, debug := os.LookupEnv("DEBUG")
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	app.Start()
	defer app.Stop()

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chanSignal
}
