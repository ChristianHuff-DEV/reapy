package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
)

var baseSuggests = []prompt.Suggest{{Text: "help", Description: "Show available commands"}, {Text: "exit", Description: "Exit the application"}, {Text: "execute", Description: "Choose a plan to execute"}}

var baseFunctions = map[string]func(){"help": funcHelp, "exit": funcExit}

var funcHelp = func() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	time.Sleep(4 * time.Second)
	s.Stop()
}

var funcExit = func() {
	os.Exit(0)
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
