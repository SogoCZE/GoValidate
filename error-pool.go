package govalidate

import "fmt"

type ErrorPoolErrors map[string]string

func (errors ErrorPoolErrors) Error() string {
	invalidFields := ""
	i := 0
	for name, _ := range errors {
		comma := ","
		if i == 0 {
			comma = ""
		}

		invalidFields += comma + name
		i += 1
	}
	return fmt.Sprintf("Tyto pole jsou nevalidní: %s", invalidFields)
}

type ErrorPool struct {
	errors ErrorPoolErrors
}

func NewErrorPool() *ErrorPool {
	return &ErrorPool{errors: make(ErrorPoolErrors)}
}

func (errorPool *ErrorPool) AddError(field string, message string) {
	errorPool.errors[field] = message
}

func (errorPool *ErrorPool) GetErrors() ErrorPoolErrors {
	return errorPool.errors
}
