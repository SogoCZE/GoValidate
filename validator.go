package govalidate

import (
	"fmt"
	"reflect"
	"strconv"
)

type StaticErrors struct {
	Object string
	Int    string
	Slice  string
	String string
	Float  string

	Required    string
	ServerError string
}

type ValidatorParam struct {
	value interface{}
}

func (validatorParam *ValidatorParam) String() string {
	stringValue, ok := validatorParam.value.(string)
	if !ok {
		panic("Cannot cast value to string")
	}

	return stringValue
}

func (validatorParam *ValidatorParam) Int() int {
	intValue, err := strconv.Atoi(validatorParam.String())
	if err != nil {
		panic("Cannot cast value to int -> " + err.Error())
	}

	return intValue
}

type ValidatorFunc = func(value reflect.Value, params []ValidatorParam) string
type ValidatorErrors = map[string][]string
type ValidatorMap map[string]ValidatorFunc

type Validator struct {
	Validators   ValidatorMap
	StaticErrors StaticErrors
}

func NewValidator() *Validator {
	return &Validator{Validators: ValidatorMap{
		"required": Required,
		"length":   Length,
		"email":    Email,
	}, StaticErrors: StaticErrors{
		Int:         "Toto pole musí obsahovat číslo",
		Object:      "Toto pole musí obsahovat objekt",
		Slice:       "Toto pole musí obsahovat pole elementů",
		String:      "Toto pole musí obsahovat text",
		Float:       "Toto pole musí obsahovat desetinné číslo",
		Required:    "Toto pole je povinné",
		ServerError: "Ou... Nastala chyba na straně serveru. Zkuste to prosím později. (%s)",
	}}
}

func (validator *Validator) SetStaticErrors(staticErrors StaticErrors) {
	validator.StaticErrors = staticErrors
}

func (validator *Validator) GenerateServerError(message string) string {
	return fmt.Sprintf(validator.StaticErrors.ServerError, message)
}

func isRequired(validators []ValidatorAST) bool {
	for _, validator := range validators {
		if validator.Name == "required" {
			return true
		}
	}

	return false
}
