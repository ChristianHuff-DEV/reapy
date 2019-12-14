package cli

import (
	"fmt"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
)

// Config represents the content of the yaml file used to define what this app is capable of doing
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

	// Does the user want to execute a plan?
	if strings.HasPrefix(command, "execute ") {
		// Extrace the name of the plan
		planName := command[len("execute "):]
		if plan, ok := Config.Plans[planName]; ok {
			executePlan(plan)
			return
		}
		fmt.Printf("Plan %s not found\n", planName)
	}

	// Find the function for the given command and execute it
	if function, ok := baseFunctions[command]; ok {
		function()
	} else {
		fmt.Println("Command not found!")
	}
}

func executePlan(plan model.Plan) {
	for _, task := range plan.Tasks {
		for _, step := range task.Steps {
			result := step.Execute()
			// Print what when wrong if an error occurred
			if !result.WasSuccessful {
				color.Red.Println(result.Message)
			}
		}
	}
}
