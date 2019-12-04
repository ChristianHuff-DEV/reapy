package model

import "github.com/c-bata/go-prompt"

// Config is the complete result of reading the yaml file and creating the possible commands available from it
type Config struct {
	Yaml Yaml
	// Completer defines which prompts are available in the terminal
	Completer func(document prompt.Document) []prompt.Suggest
	// Executor is the mapping between a command and the function to run
	Executor func(command string)
}

// Yaml represents the root structure of the yaml file defining the plans
type Yaml struct {
	Variables map[string]string
	Plans     []Plan
}

// Plan defines the structure of the equally named node in the config yaml
type Plan struct {
	Name  string
	Tasks []Task
}

// Task defines the structure of the equally named node in the config yaml
type Task struct {
	Name  string
	Steps []Step
}

// Step is the interface defining the methods each step has to implement to be generically executable
type Step interface {
	// GetKind returns the type this step is (i.e. Download,Unzip etc)
	GetKind() string
	// GetDescription return description of what the step does
	GetDescription() string
	// Execute performs the actual work the step is responsible for
	Execute() Result
	// FromConfig read the given map (representing the yaml definition) and creates a step instance of it
	FromConfig(configYaml map[string]interface{})
}

// RunnableStep defines the base struct
type RunnableStep struct {
	Kind        string
	Description string
}

// Result is the outcome of executing a task.
type Result struct {
	Message       string
	WasSuccessful bool
}
