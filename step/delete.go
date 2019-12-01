package step

import (
	"os"

	"github.com/ChristianHuff-DEV/reapy/model"
)

type Delete struct {
	model.RunnableStep
	Path string
}

func (this Delete) GetKind() string {
	return this.Kind
}

func (this Delete) GetDescription() string {
	return this.Description
}

func (this Delete) Execute() model.Result {
	return delete(this.Path)
}

func delete(path string) (result model.Result) {
	f, err := os.Stat(path)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	if f.IsDir() {
		err := os.RemoveAll(path)
		if err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
	} else {
		err := os.Remove(path)
		if err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
	}

	result.WasSuccessful = true
	result.Message = path + " deleted"
	return result
}
