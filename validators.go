package govalidate

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
)

func Required(rawValue reflect.Value, params []ValidatorParam) string {

	if rawValue.Type().Name() == "string" {
		value := rawValue.String()

		if len(value) < 1 || value == "null" {
			return "Vyplňte prosím toto pole."
		}

		return ""
	}

	if rawValue.Type().Name() == "int" {
		// TODO: How to handle this (default struct int value is zero and user value can be also zero and both can be technically valid)
		// Zero is not dangerous value so let this go. (we can fix this with Length validator)
		log.Printf("[Validator - Required] warning: int field")
		return ""
	}

	log.Printf("[Validator - Required] unhadled type!")

	return ""
}

func Length(rawValue reflect.Value, params []ValidatorParam) string {
	min := params[0].Int()
	max := params[1].Int()

	errorMessage := fmt.Sprintf("Toto pole vyžaduje rozsah délky %d-%d znaků.", min, max)
	valueType := rawValue.Type().Name()

	if valueType == "string" {
		valueLen := len(rawValue.String())

		if valueLen < min || valueLen > max {
			return errorMessage
		}

		return ""
	}

	if valueType == "int" {
		value := int(rawValue.Int())

		if value < min || value > max {
			return errorMessage
		}

		return ""
	}

	log.Printf("[Validator - Length] unhadled type!")

	return ""
}

func Email(rawValue reflect.Value, params []ValidatorParam) string {
	errorVal := "Zadejte validní email"
	value := rawValue.String()

	matched, err := regexp.Match("^[a-z0-9!#$%&'*+\\/=?^_`{|}~.-]+@[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$", []byte(value))
	if err != nil {
		log.Println(err)
		return errorVal
	}

	if !matched {
		return errorVal
	}

	return ""
}
