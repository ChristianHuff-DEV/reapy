package step

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindCommand defines the name for a command step in the config file
const KindCommand = "Command"

// Command executes the defined command
type Command struct {
	model.RunnableStep
	Command string
	Args    []string
	Path    string
}

// GetKind returns the kind this step represents
func (command Command) GetKind() string {
	return command.Kind
}

// GetDescription gives a description of what this step does
func (command Command) GetDescription() string {
	return command.Description
}

// FromConfig create a command struct from the given config
func (command *Command) FromConfig(stepConfig map[string]interface{}) {
	command.Kind = KindCommand
	preferencesYaml := stepConfig["Preferences"].(map[string]interface{})
	command.Command = preferencesYaml["Command"].(string)
	var args []string
	for _, arg := range preferencesYaml["Args"].([]interface{}) {
		args = append(args, arg.(string))
	}
	if path, ok := preferencesYaml["Path"].(string); ok {
		command.Path = path
	} else {
		command.Path = ""
	}
}

// Execute runs the defined command
func (command Command) Execute() (result model.Result) {
	log.Printf("Executing: %s in %s", command.Command, command.Path)
	// Create the command
	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.Path

	var stdBuffer bytes.Buffer

	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Print(err)
		return model.Result{
			Message:       err.Error(),
			WasSuccessful: false,
		}
	}

	// Output of the comman executed. Could later be used to have some logic running on it to
	// determine if the execution was successful.
	//output := stdBuffer.String()

	result.WasSuccessful = true
	result.Message = "The command \"" + command.Command + "\"" + " was executed successfully."

	return result
}
