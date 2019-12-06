package config

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
)

var baseSuggests = []prompt.Suggest{{Text: "help", Description: "Show available commands"}, {Text: "exit", Description: "Exit the application"}}

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
