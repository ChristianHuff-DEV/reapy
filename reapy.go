//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/ChristianHuff-DEV/reapy/cli"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	file, err := setupLogging()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cli.Config, err = readPlanDefinition()
	if err != nil {
		color.Red.Printf("Error reading plans definition file: %s\n", err)
		log.Fatal(err)
	}

	p := prompt.New(cli.Executor, cli.Completer)
	p.Run()
}

func setupLogging() (*os.File, error) {
	log.SetFormatter(&prefixed.TextFormatter{
            DisableColors: true,
            TimestampFormat : "2006-01-02 15:04:05",
            FullTimestamp:true,
						ForceFormatting: true,
        },)

	// Create log file
	file, err := os.OpenFile("reapy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	// Set file as log output
	return file, nil
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() (config model.Config, err error) {
	log.Println("reading plans configuration file")
	config, err = cli.Extract("test.yaml")
	if err != nil {
		return config, err
	}
	return config, nil
}
