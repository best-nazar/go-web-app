package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Errors(c *gin.Context) map[string]string {
	var data = make(map[string]string)

	for _, err := range c.Errors {
		switch err.Err.(type) {
			case validator.ValidationErrors:
				formErrors(err, data)
			case error:
				convertErrorTypeString(err, data)
			default:
				convertDefault(err)
		}
	}

	return data
}

// error data will appear below the input field of a Form
func formErrors(err *gin.Error, data map[string]string) {
	e:=err.Err.(validator.ValidationErrors)

	for _, v := range e {
		data[v.StructField()] = fmt.Sprintf("Value '%s' failed the validation because of constraint '%s' %s",v.Value().(string), v.Tag(), v.Param())
	}
}

// "massage" key data will appear on header error alert
func convertErrorTypeString(err error, data map[string]string) {
	data["message"] = Capitalize(err.Error())
}

func convertDefault(err *gin.Error) {
	panic("NOT implemented 'convertDefault' in prettyErrors.go on " + err.Err.Error())
}
