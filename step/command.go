package step

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindCommand defines the name for a command step in the config file
const KindCommand = "Command"

const fieldNamePreferences = "Preferences"
const fieldNameCommand = "Command"
const fieldNamePath = "Path"
const fieldNameArgs = "Args"
const fieldNameSilent = "Silent"

// Command executes the defined command
type Command struct {
	model.RunnableStep
	Command string
	Args    []string
	Path    string
	// Whether or not the output of the command is printed to the console
	Silent bool
}

// GetKind returns the kind this step represents
func (command Command) GetKind() string {
	return command.Kind
}

// GetDescription gives a description of what this step does
func (command Command) GetDescription() string {
	return command.Description
}

// FromConfig creates a command struct from the given config
func (command *Command) FromConfig(stepConfig map[string]interface{}) error {
	command.Kind = KindCommand
	if description, ok := stepConfig["Description"]; ok {
		command.Description = description.(string)
	}
	preferencesYaml := stepConfig[fieldNamePreferences].(map[string]interface{})
	command.Command = preferencesYaml[fieldNameCommand].(string)

	// Extract the "Args" field if the exist
	if argsYaml, ok := preferencesYaml[fieldNameArgs].([]interface{}); ok {
		var args []string
		for _, arg := range argsYaml {
			args = append(args, arg.(string))
		}
		command.Args = args
	}

	// Extract "Path" field
	if path, ok := preferencesYaml[fieldNamePath].(string); ok {
		command.Path = path
	} else {
		command.Path = ""
	}

	// Extract "Silent" field
	if silent, ok := preferencesYaml[fieldNameSilent].(bool); ok {
		command.Silent = silent
	} else {
		silent = false
	}

	return nil
}

// Execute runs the defined command
func (command Command) Execute() (result model.Result) {
	fmt.Println(command.Description)
	log.Printf("Executing: %s in %s", command.Command, command.Path)
	// Create the command
	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.Path

	var writer io.Writer
	var stdBuffer bytes.Buffer

	// If it's silent we only write to the buffer
	if command.Silent {
		writer = io.MultiWriter(&stdBuffer)
	} else {
		writer = io.MultiWriter(os.Stdout, &stdBuffer)
	}

	cmd.Stdout = writer
	cmd.Stderr = writer

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Print(err)
		return model.Result{
			Message:       err.Error(),
			WasSuccessful: false,
		}
	}

	// Read the output of the executed command
	// For now we just log it. But it could be later used to determine the success of running the command.
	log.Println("Command output:")
	log.Println("-----------------------------------")
	scanner := bufio.NewScanner(&stdBuffer)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
	log.Println("-----------------------------------")

	result.WasSuccessful = true
	result.Message = "The command \"" + command.Command + "\"" + " was executed successfully."

	return result
}
