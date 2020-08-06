package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	err := Execute()
	if err != nil {
		log.Info(err)
		os.Exit(1)
	}
	log.Debug("Scan Complete")
}
