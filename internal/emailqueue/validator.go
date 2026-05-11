package emailqueue

import (
	"fmt"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`)

func validateEmail(email string) error {
	switch {
	case email == "":
		return fmt.Errorf("email is empty")
	case !emailRegex.MatchString(email):
		return fmt.Errorf("email format is invalid")
	default:
		return nil
	}
}
