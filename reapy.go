//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianHuff-DEV/reapy/cli"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	file, err := setupLogging()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = cli.InitializePlans()
	if err != nil {
		log.Fatal(err)
	}

	cli.Start()
}

func setupLogging() (*os.File, error) {
	log.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	// Create log file
	file, err := os.OpenFile("reapy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	// Set file as log output
	return file, nil
}
