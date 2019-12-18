package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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
			sort.Slice(suggests, func(i, j int) bool {
				x := strings.Compare(suggests[i].Text, suggests[j].Text)
				return x < 0
			})
		}
		return prompt.FilterFuzzy(suggests, command, true)
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
	// Label used to break out of the nested loop if a step fails and the user chooses not to continue
out:
	for _, task := range plan.Tasks {
		for _, step := range task.Steps {
			result := step.Execute()
			// Handle a step failing
			if !result.WasSuccessful {
				color.Red.Println(result.Message)
				response := "Abort, Continue or Retry?"
				prompt := &survey.Select{
					Message: "",
					Options: []string{"Abort", "Continue"},
				}
				survey.AskOne(prompt, &response)
				switch response {
				case "Abort":
					// Break out of th nested loop
					break out
				case "Continue":
					// Do nothing. Just continue with the execution.
				}
			}
		}
	}
}
