package helpers

import "errors"

func ValidateUserPassword (pass1, pass2 string) error {
	if pass1 != pass2 {
		return errors.New("passwords do not match")
	}

	return nil
}