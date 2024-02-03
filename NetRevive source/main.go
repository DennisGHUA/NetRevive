package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	LogInfo("Running NetRevive")

	// Block the main thread until a signal is received
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer LogInfo("Exiting NetRevive")
	<-quit
}
