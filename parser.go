package govalidate

import (
	"strings"
)

type ValidatorAST struct {
	Name   string
	Params []string
}

func eatToThe(str *string, ends []rune) (string, rune) {
	charsToEat := 0
	eaten := ""
	end := rune(0)
	b := *str

	for _, char := range *str {
		charsToEat += 1
		foundEnd := false

		for _, endRune := range ends {
			if endRune == char {
				end = endRune
				foundEnd = true
			}
		}

		if foundEnd {
			break
		}

		eaten += string(char)
	}

	*str = b[charsToEat:]
	return eaten, end
}

func getLastIndex(validators *[]ValidatorAST) int {
	length := len(*validators)

	if length == 0 {
		return 0
	}

	return length - 1
}

// TODO: make some error boundaries
func parse(validationTag string) []ValidatorAST {
	insideValidatorParams := false
	insideString := false
	createNewValidator := true
	validators := []ValidatorAST{}

	for len(validationTag) != 0 {
		var currentPossibleEnds []rune

		if insideString {
			currentPossibleEnds = []rune{'\''}
		} else if insideValidatorParams {
			currentPossibleEnds = []rune{')', ',', '\''}
		} else {
			currentPossibleEnds = []rune{' ', '('}
		}

		value, endRune := eatToThe(&validationTag, currentPossibleEnds)

		if createNewValidator && len(value) > 0 {
			validators = append(validators, ValidatorAST{Name: value, Params: []string{}})
			createNewValidator = false
		}

		if endRune == '\'' {

			// This is the closing one
			if insideString {
				validators[getLastIndex(&validators)].Params = append(validators[getLastIndex(&validators)].Params, value)
				insideString = false
			} else {
				insideString = true
			}

		}

		if endRune == ',' && len(value) > 0 {
			param := strings.ReplaceAll(value, " ", "")
			validators[getLastIndex(&validators)].Params = append(validators[getLastIndex(&validators)].Params, param)
		}

		if endRune == ' ' {
			createNewValidator = true
		}

		if endRune == '(' {
			insideValidatorParams = true
		}

		if endRune == ')' {
			insideValidatorParams = false

			if len(value) > 0 {
				param := strings.ReplaceAll(value, " ", "")
				validators[getLastIndex(&validators)].Params = append(validators[getLastIndex(&validators)].Params, param)
			}

			createNewValidator = true
		}

	}

	return validators
}
