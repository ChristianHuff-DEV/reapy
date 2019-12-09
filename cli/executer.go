package cli

import (
	"log"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
)

// Execute runs all steps in all tasks of the given plan
func Execute(steps model.Plan) {
	for _, task := range steps.Tasks {
		for _, step := range task.Steps {
			result := step.Execute()
			// Print what when wrong if an error occurred
			if !result.WasSuccessful {
				log.Print(result.Message)
			}
		}
	}
}

// Completer determines the suggestions shown to the user
var Completer = func(document prompt.Document) (suggests []prompt.Suggest) {
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

// Executor determines which what to do with the given command
var Executor = func(command string) {

	// Find the function for the given command and execute it
	if function, ok := baseFunctions[command]; ok {
		function()
	}
}

// HasSubcommands returns false for all commands that don't have a subcommand
func hasSubcommands(word string) bool {
	return !(strings.HasPrefix(word, "help") || strings.HasPrefix(word, "exit"))
}
