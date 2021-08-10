package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func SchemaValidator(schema interface{}) []string {
	validate := validator.New()
	err := validate.Struct(schema)
	var error []string
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			// return
		}

		for _, err := range err.(validator.ValidationErrors) {

			error = append(error, err.Field()+" "+err.Tag()+" "+err.Param())
		}
	}
	return error
}
