package helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BodyError struct {
	Context string
	Message string
}

var prettyErrors []*BodyError
var data *BodyError

func Errors(c *gin.Context) []*BodyError {
	prettyErrors = nil // reset error list

	for _, err := range c.Errors {
		switch err.Err.(type) {
		case validator.ValidationErrors:
			convertErrorTypePrivate(err)
		case error:
			convertErrorTypeString(err)
		default:
			convertDefault(err)
		}
	}

	return prettyErrors
}

func convertErrorTypePrivate(err *gin.Error) {
	e:=err.Err.(validator.ValidationErrors)
	
	for _, v := range e {
		data = &BodyError{
			Context: Capitalize(v.StructField() + ":"),
			Message: fmt.Sprintf("Value '%s' failed the validation because of constraint '%s'", v.Value().(string), v.Tag()),
		}

		prettyErrors = append(prettyErrors, data)
	}
}

func convertErrorTypeString(err error) {
		msgs := strings.Split(err.Error(), "|")

		if len(msgs) == 2 {
			data = &BodyError{
				Context: Capitalize(msgs[0] + ":"),
				Message: msgs[1],
			}

		} else {
			data = &BodyError{
				Context: "Error:",
				Message: msgs[0],
			}
		}
	prettyErrors = append(prettyErrors, data)
}

func convertDefault(err *gin.Error) {
	panic("NOT implemented 'convertDefault' in prettyErrors.go on " + err.Err.Error())
}