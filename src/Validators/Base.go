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

var validate *validator.Validate

func Validate(data interface{}) (bool, Errors) {
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

func validateVariable() {

	myEmail := "joeybloggs.gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}

	// email ok, move on
}
