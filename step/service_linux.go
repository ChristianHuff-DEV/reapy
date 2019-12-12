package step

import (
	"fmt"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// FromConfig creates the struct representation of a service step
func (service *Service) FromConfig(configYaml map[string]interface{}) error {
	return fmt.Errorf("the step type \"%s\" is not supported on this operating system", KindService)
}

// Execute starts/stops a service
func (service *Service) Execute() model.Result {
	return model.Result{WasSuccessful: false, Message: "the step type is not supported on this operating system"}
}
