package govalidate

import (
	"fmt"
	"reflect"
)

func (validator *Validator) validateStructSlice(kind string, name string, fieldValue reflect.Value, errorPool *ErrorPool, validators []ValidatorAST) bool {
	sliceType := fieldValue.Type().Elem().Kind()

	for index := 0; index < fieldValue.Len(); index++ {

		if sliceType == reflect.Struct {
			validator.validateStruct(fmt.Sprintf("%s[%d]", name, index), errorPool, fieldValue.Index(index))
		} else if sliceType == reflect.Slice {
			itemName := fmt.Sprintf("%s[%d]", name, index)

			skipToNext := validator.validateStructSlice(sliceType.String(), itemName, fieldValue.Index(index), errorPool, validators)
			if skipToNext {
				continue
			}

			for _, fieldValidator := range validators {

				var params []ValidatorParam
				for _, param := range fieldValidator.Params {
					params = append(params, ValidatorParam{value: param})
				}

				validatorError := validator.Validators[fieldValidator.Name](fieldValue.Index(index), params)
				if len(validatorError) > 0 {
					errorPool.AddError(itemName, validatorError)
				}

			}
		} else {
			itemName := fmt.Sprintf("%s[%d]", name, index)

			for _, fieldValidator := range validators {

				var params []ValidatorParam
				for _, param := range fieldValidator.Params {
					params = append(params, ValidatorParam{value: param})
				}

				validatorError := validator.Validators[fieldValidator.Name](fieldValue.Index(index), params)
				if len(validatorError) > 0 {
					errorPool.AddError(itemName, validatorError)
				}

			}
		}

	}

	return false
}

func (validator *Validator) validateStruct(parent string, errorPool *ErrorPool, levelValue reflect.Value) {
	levelType := levelValue.Type()

	if len(parent) > 0 {
		parent = parent + "."
	}

	for i := 0; i < levelType.NumField(); i++ {
		field := levelType.Field(i)
		fieldValue := levelValue.Field(i)
		kind := field.Type.Kind().String()

		validationTag := field.Tag.Get("validation")
		validators := parse(validationTag)
		jsonName := field.Tag.Get("json")

		if len(validationTag) == 0 {
			continue
		}

		if kind == "struct" {
			validator.validateStruct(parent+jsonName, errorPool, fieldValue)
		} else if kind == "slice" {
			skipToNext := validator.validateStructSlice(kind, parent+jsonName, fieldValue, errorPool, validators)
			if skipToNext {
				continue
			}
		}

		// Handle validators for current field.
		for _, fieldValidator := range validators {

			var params []ValidatorParam
			for _, param := range fieldValidator.Params {
				params = append(params, ValidatorParam{value: param})
			}

			validatorError := validator.Validators[fieldValidator.Name](fieldValue, params)
			if len(validatorError) > 0 {
				errorPool.AddError(parent+jsonName, validatorError)
			}

		}

	}

}

// Validates Struct.
func (validator *Validator) Validate(data interface{}, errorPool *ErrorPool) bool {
	validator.validateStruct("", errorPool, reflect.ValueOf(data).Elem())

	return !(len(errorPool.GetErrors()) > 0)
}
