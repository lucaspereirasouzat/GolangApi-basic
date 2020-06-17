package validatores

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

type Error struct {
	Field   string
	Message string
}

type Errors struct {
	Errors []Error
}

func Validate(data interface{}) (bool, Errors) {
	var validate *validator.Validate
	err := validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			//panic("PROBLEMA COM A VALIDAÇÃO")
			return false, Errors{}
		}
		var listError Errors
		for i, err := range err.(validator.ValidationErrors) {
			err := Error{
				err.Field(),
				"erro mensagem" + err.Tag(),
			}

			listError.Errors[i] = err
		}
		// from here you can create your own error messages in whatever language you wish
		return true, listError
	}
	return false, Errors{}
}
