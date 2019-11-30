package step

import (
	"log"

	"github.com/ChristianHuff-DEV/reapy/model"
)

type Unzip struct {
	model.RunnableStep
}

func (this Unzip) GetKind() string {
	return this.Kind
}

func (this Unzip) GetDescription() string {
	return this.Description
}

func (this Unzip) Execute() {
	log.Print("Download Unzip()")
}
