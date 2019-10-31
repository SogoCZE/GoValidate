package govalidate

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func (validator *Validator) validateJsonSlice(kind string, name string, fieldValue reflect.Value, data *json.RawMessage, errorPool *ErrorPool, validators []ValidatorAST) bool {
	var value []*json.RawMessage
	err := json.Unmarshal(*data, &value)
	if err != nil {
		errorPool.AddError(name, validator.StaticErrors.Slice)
		return true
	}

	sliceType := fieldValue.Type().Elem().Kind().String()
	fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), len(value), len(value)))

	if sliceType == "struct" {
		for index, item := range value {
			var jsonDataSlice map[string]*json.RawMessage
			err := json.Unmarshal(*item, &jsonDataSlice)
			if err != nil {
				errorPool.AddError(name, validator.StaticErrors.Object)
				continue
			}
			validator.validateJsonStruct(fmt.Sprintf("%s[%d]", name, index), errorPool, fieldValue.Index(index), jsonDataSlice)
		}
	} else if sliceType == "splice" {
		for index, item := range value {
			itemName := fmt.Sprintf("%s[%d]", name, index)

			skipToNext := validator.validateJsonSlice(sliceType, itemName, fieldValue.Index(index), item, errorPool, validators)
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
		}
	} else {
		for index, item := range value {
			itemName := fmt.Sprintf("%s[%d]", name, index)

			skipToNext := validator.validateJsonPrimitive(sliceType, itemName, fieldValue.Index(index), item, errorPool)
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
		}
	}

	return false
}

func (validator *Validator) validateJsonPrimitive(kind string, name string, fieldValue reflect.Value, data *json.RawMessage, errorPool *ErrorPool) bool {
	if kind == "string" {
		var value string
		err := json.Unmarshal(*data, &value)
		if err != nil {
			errorPool.AddError(name, validator.StaticErrors.String)
			return true
		}

		fieldValue.SetString(value)
	}

	if kind == "int" {
		var value int64
		err := json.Unmarshal(*data, &value)
		if err != nil {
			errorPool.AddError(name, validator.StaticErrors.Int)
			return true
		}

		fieldValue.SetInt(value)
	}

	if kind == "float" {
		var value float64
		err := json.Unmarshal(*data, &value)
		if err != nil {
			errorPool.AddError(name, validator.StaticErrors.Float)
			return true
		}

		fieldValue.SetFloat(value)
	}

	return false
}

func (validator *Validator) validateJsonStruct(parent string, errorPool *ErrorPool, levelValue reflect.Value, inJsonData map[string]*json.RawMessage) {
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

		currentDataInJson, ok := inJsonData[jsonName]
		if !ok {
			if isRequired(validators) {
				errorPool.AddError(parent+jsonName, validator.StaticErrors.Required)
			}
			continue
		}

		if kind == "struct" {
			var jsonData map[string]*json.RawMessage
			err := json.Unmarshal(*currentDataInJson, &jsonData)
			if err != nil {
				errorPool.AddError(parent+jsonName, validator.StaticErrors.Object)
				continue
			}

			validator.validateJsonStruct(parent+jsonName, errorPool, fieldValue, jsonData)
		} else if kind == "slice" {
			skipToNext := validator.validateJsonSlice(kind, parent+jsonName, fieldValue, currentDataInJson, errorPool, validators)
			if skipToNext {
				continue
			}
		} else {
			skipToNext := validator.validateJsonPrimitive(kind, parent+jsonName, fieldValue, currentDataInJson, errorPool)
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

// Validates JSON and bind data intro it.
func (validator *Validator) ValidateJson(rawData []byte, bindTo interface{}, errorPool *ErrorPool) bool {
	if !json.Valid(rawData) {
		errorPool.AddError("other", validator.GenerateServerError("invalid json encoding"))
		return false
	}

	jsonData := make(map[string]*json.RawMessage)
	err := json.Unmarshal(rawData, &jsonData)
	if err != nil {
		log.Printf("[Validator] Base level unmarshal error error -> %s", err)
		errorPool.AddError("other", validator.GenerateServerError("base level unmarshal error"))
		return false
	}

	validator.validateJsonStruct("", errorPool, reflect.ValueOf(bindTo).Elem(), jsonData)

	return !(len(errorPool.GetErrors()) > 0)
}
