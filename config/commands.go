package config

import (
	"os"
	"strings"
	"time"

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

var completer = func(document prompt.Document) (suggests []prompt.Suggest) {
	command := document.Text

	// Return empty prompts if the given command has no sub-commands
	if hasSubcommands(command) {
		return []prompt.Suggest{}
	}

	// Add all plans for the execute command
	if strings.HasPrefix(command, "execute") {

	}

	return prompt.FilterHasPrefix(baseSuggests, command, true)
}

// HasSubcommands returns false for all commands that don't have a subcommand
func hasSubcommands(word string) bool {
	return !(strings.HasPrefix(word, "help") || strings.HasPrefix(word, "exit"))
}

var executor = func(command string) {

	// Find the function for the given command and execute it
	if function, ok := baseFunctions[command]; ok {
		function()
	}
}
