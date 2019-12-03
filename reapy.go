//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"log"
	"os"

	c "github.com/ChristianHuff-DEV/reapy/config"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
)

var config model.Config

func main() {
	// for _, plan := range config.Plans {
	// 	c.Execute(plan)
	// }
	p := prompt.New(executor, completer)
	p.Run()
}

// completer provides the available commands
func completer(document prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "help"},
		{Text: "execute"},
		{Text: "exit"},
	}
}

func executor(command string) {
	switch command {
	case "help":
		log.Print("Executed help")
	case "execute":
		log.Print("Execute execute")
	case "exit":
		os.Exit(0)
	}

}

// init will read the config yaml before starting the app itself
func init() {
	config = readPlanDefinition()
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() model.Config {
	return c.Extract("test.yaml")
}
