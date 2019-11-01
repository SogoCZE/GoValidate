package govalidate

import "testing"

type BasicJson struct {
	Name string `json:"name" validation:"required length(0, 42)"`
	Age  int    `json:"age" validation:"required"`
}

func TestBasicJsonCorrect(t *testing.T) {
	validator := NewValidator()
	errorPool := NewErrorPool()

	dataStruct := BasicJson{}
	data := []byte(`
		{
			"name": "John Doe",
			"age": 38			
		}
	`)

	result := validator.ValidateJson(data, &dataStruct, errorPool)
	if !result {
		t.Fatal(errorPool.GetErrors())
		return
	}

	if dataStruct.Name != "John Doe" || dataStruct.Age != 38 {
		t.Fatal(dataStruct)
		return
	}
}

func TestBasicJsonIncorrect(t *testing.T) {
	validator := NewValidator()
	errorPool := NewErrorPool()

	dataStruct := BasicJson{}
	data := []byte(`
		{
			"name": "John Doe",
		}
	`)

	resultCorrect := validator.Validate(data, &dataStruct, errorPool)
	if resultCorrect {
		t.Fatal(dataStruct)
		return
	}
}
