package step

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/ChristianHuff-DEV/reapy/model"
)

type Command struct {
	model.RunnableStep
	Command string
	Args    []string
	Path    string
}

func (this Command) GetKind() string {
	return this.Kind
}

func (this Command) GetDescription() string {
	return this.Description
}

func (this Command) Execute() (result model.Result) {
	log.Printf("Executing: %s", this.Command)
	// Create the command
	cmd := exec.Command(this.Command, this.Args...)
	cmd.Dir = this.Path

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
	result.Message = "The command \"" + this.Command + "\"" + " was executed successfully."

	return result
}
