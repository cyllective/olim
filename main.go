package main

import (
	"os"
	"os/signal"
	"syscall"

	logger "github.com/rtfmkiesel/kisslog"

	"github.com/cyllective/olim/internal/app"
)

var version = "@DEV"

func main() {
	if err := logger.InitDefault("github.com/cyllective/olim" + version); err != nil {
		panic(err)
	}

	_, debug := os.LookupEnv("DEBUG")
	logger.FlagDebug = debug

	app.Start()
	defer app.Stop()

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chanSignal
}
