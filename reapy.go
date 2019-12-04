//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	c "github.com/ChristianHuff-DEV/reapy/config"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
)

var config model.Config

func main() {
	// for _, plan := range config.Plans {
	// 	c.Execute(plan)
	// }
	p := prompt.New(executor, completer)
	p.Run()
}

// completer provides the available commands
func completer(document prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "help"},
		{Text: "execute"},
		{Text: "exit"},
	}
}

func executor(command string) {
	switch command {
	case "help":
		showWait()
		log.Print("Executed help")
		showChecklist()
	case "execute":
		log.Print("Execute execute")
		askQuestion()
	case "exit":
		os.Exit(0)
	}

}

func askQuestion() {
	var qs = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is your name?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		}}
	answers := struct {
		Name string // survey will match the question and field names
	}{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("chose %s.\n", answers.Name)
}

func showWait() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                    // Start the spinner
	time.Sleep(4 * time.Second)                                  // Run for some time to simulate work
	s.Stop()
}

func showChecklist() {
	days := []string{}
	prompt := &survey.MultiSelect{
		Message: "What days do you prefer:",
		Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	}
	survey.AskOne(prompt, &days)
}

// init will read the config yaml before starting the app itself
func init() {
	config = readPlanDefinition()
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() model.Config {
	return c.Extract("test.yaml")
}
