package cli

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
)

var p *prompt.Prompt

var baseSuggests = []prompt.Suggest{{Text: "help", Description: "Show available commands"}, {Text: "exit", Description: "Exit the application"}, {Text: "execute", Description: "Choose a plan to execute"}, {Text: "reload", Description: "Reload plans from config"}}

var baseFunctions = map[string]func(){"help": funcHelp, "exit": funcExit, "reload": funcReload}

var funcHelp = func() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	time.Sleep(4 * time.Second)
	s.Stop()
}

var funcExit = func() {
	os.Exit(0)
}

var funcReload = func() {
	fmt.Println("funcReload")
	InitializePlans()
}

// Start creates the prompt instance and runs it
func Start() {
	p = prompt.New(Executor, Completer)
	p.Run()
}

// InitializePlans creates the global config instance that creates the plans definitions
//
// Can also be used to update the available plans.
func InitializePlans() (err error) {
	Config, err = readPlanDefinition()
	if err != nil {
		color.Red.Printf("Error reading plans definition file: %s\n", err)
		log.Fatal(err)
		return err
	}
	return nil
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() (config model.Config, err error) {
	log.Println("reading plans configuration file")
	config, err = Extract("test.yaml")
	if err != nil {
		return config, err
	}
	return config, nil
}
