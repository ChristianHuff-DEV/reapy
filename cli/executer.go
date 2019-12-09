package cli

import (
	"fmt"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
)

var Config model.Config

// Completer determines the suggestions shown to the user
var Completer = func(document prompt.Document) (suggests []prompt.Suggest) {
	// The current command (everything between the beginning of the line and the next space or between two spaces)
	command := document.GetWordBeforeCursor()
	text := document.Text

	// If the command is "execute " show the available plans
	if strings.HasPrefix(text, "execute ") {
		for _, plan := range Config.Plans {
			suggests = append(suggests, prompt.Suggest{Text: plan.Name, Description: plan.Description})
		}
		return suggests
	}

	return prompt.FilterHasPrefix(baseSuggests, command, true)
}

// Executor determines which what to do with the given command
var Executor = func(command string) {
	// Find the function for the given command and execute it
	if function, ok := baseFunctions[command]; ok {
		function()
	} else {
		fmt.Println("Command not found!")
	}
}
