package step

import (
	"fmt"

	"github.com/ChristianHuff-DEV/reapy/model"
)

const KindService = "Service"

type Service struct {
	model.RunnableStep
	ServiceName string
	Action      string
}

func (service Service) GetKind() string {
	return service.Kind
}

func (service Service) GetDescription() string {
	return service.Description
}

func (service *Service) FromConfig(configYaml map[string]interface{}) error {
	return fmt.Errorf("the step type \"%s\" is not supported on this operating system", KindService)
}

func (service *Service) Execute() model.Result {
	return model.Result{WasSuccessful: false, Message: "the step type is not supported on this operating system"}
}
