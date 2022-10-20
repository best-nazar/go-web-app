// Error handler based on the gin binding tag
package errorSrc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func MakeErrorView(title string, err error) ErrorView {
	errView := ErrorView{
		Title: title, 
		Messages: parseError(err),
	}
	
	return errView
}

func MakeErrorViewFrom(title, field string, code int) ErrorView {
	errView := ErrorView{
		Title: title, 
		Messages: codeError(field, code),
	}
	
	return errView
}

// Parsing array of Tag errors.
// parseError takes an error or multiple errors and attempts to determine the best path to convert them into
// human readable strings
func parseError(errs ...error) []string {
	var out []string
	for _, err := range errs {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			// if the type is validator.ValidationErrors then it's actually an array of validator.FieldError so we'll
			// loop through each of those and convert them one by one
			for _, e := range typedError {
				out = append(out, parseFieldError(e))
			}
		case *json.UnmarshalTypeError:
			// similarly, if the error is an unmarshalling error we'll parse it into another, more readable string format
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}
	return out
}

// Add case with the tag for handling the message
func parseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	fieldPrefix := fmt.Sprintf("The field '%s'", e.Field())
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required_without":
		return fmt.Sprintf("%s is required if '%s' is not supplied", fieldPrefix, e.Param())
	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s symb.", fieldPrefix, param)
	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s symb.", fieldPrefix, param)
	case "min":
		return fmt.Sprintf("%s must be greater than %s symb.", fieldPrefix, e.Param())
	case "alphanum":
		return fmt.Sprintf("%s allows only Alphanumeric symb.", fieldPrefix)
	case "excludesall":
		return fmt.Sprintf("%s contains redundant '%s' symb.", fieldPrefix, e.Param())
	default:
		// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return e.Translate(translatorInstance)
		} else {
			return fmt.Errorf("%v", e).Error()
		}
	}
}

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}

func codeError(field string, code int) []string {
	var out []string
	switch code {
		case http.StatusNotFound:
			out = append(out, fmt.Sprintf("The field %s is invalid or object is not found", field))
	}

	return out
}