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

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			//panic("PROBLEMA COM A VALIDAÇÃO")
			return false, Errors{}
		}
		var listError Errors
		for i, err := range err.(validator.ValidationErrors) {

			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
			err := Error{
				err.Field(),
				"erro mensagem" + err.Tag(),
			}

			listError.Errors[i] = err

		}
		//fmt.Println(listError.Errors[0].Message)
		//c.JSON(400, listError)
		// from here you can create your own error messages in whatever language you wish
		return true, listError
	}
	return false, Errors{}
}
