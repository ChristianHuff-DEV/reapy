package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
)

var p *prompt.Prompt

var baseSuggests = []prompt.Suggest{{Text: "help", Description: "Show available commands"}, {Text: "exit", Description: "Exit the application"}, {Text: "execute", Description: "Choose a plan to execute"}, {Text: "reload", Description: "Reload plans from config"}}

var baseFunctions = map[string]func(){"help": help, "exit": exit, "reload": reload}

func help() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	time.Sleep(4 * time.Second)
	s.Stop()
}

func exit() {
	os.Exit(0)
}

// reloadPlans reads all plans again und update the suggestions of the prompt.
//
// If reading the plans returns an error the user can choose to retry it or exit the programm. If the user chooses to retry this method is called recursively.
func reload() {
	err := InitializePlans()
	if err != nil {
		response := true
		prompt := &survey.Confirm{
			Default: true,
			Message: "Retry?",
		}
		survey.AskOne(prompt, &response)
		if response {
			reload()
		} else {
			os.Exit(0)
		}
	}
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
		return err
	}
	return nil
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() (config model.Config, err error) {
	findPlanDefinitionFiles()
	log.Println("reading plans configuration file")
	config, err = Extract("test.yaml")
	if err != nil {
		return config, err
	}
	return config, nil
}

func findPlanDefinitionFiles() (files []string, err error) {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return files, err
	}
	fmt.Printf("Current directory: %s\f", dir)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Skip folders
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		fmt.Println(filepath.Ext(path))

		return nil
	})
	if err != nil {
		panic(err)
	}

	return files, nil
}
