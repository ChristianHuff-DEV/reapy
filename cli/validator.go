package cli

import (
	"fmt"
	"regexp"
)

func validateVariables(variablesYaml map[string]interface{}) error {

	for key := range variablesYaml {
		err := validateVariableName(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateVariableName(name string) error {

	regex, _ := regexp.Compile("^[_A-z0-9]*((-|_)*[_A-z0-9])*$")
	isValid := regex.MatchString(name)

	if !isValid {
		return fmt.Errorf("the variable name \"%s\" is invalid. Variable names can only consist out of characters, numbers, \"-\" and \"_\"", name)
	}

	return nil
}
