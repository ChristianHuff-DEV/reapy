package step

import "github.com/ChristianHuff-DEV/reapy/model"

import "os"

import "log"

type CreateFolder struct {
	model.RunnableStep
	Path string
}

func (this CreateFolder) GetKind() string {
	return this.Kind
}

func (this CreateFolder) GetDescription() string {
	return this.Description
}

func (this CreateFolder) Execute() model.Result {
	return create(this.Path)
}

func create(path string) model.Result {
	log.Printf("Create folder %s", path)
	//Attempt to create the directory and ignore any issues
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		log.Print(err)
	}
	return model.Result{WasSuccessful: true}
}
