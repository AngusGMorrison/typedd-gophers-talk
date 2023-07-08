package users

import "fmt"

// InvalidUserError is returned when user inputs fail validation.
type InvalidUserError struct {
	cause error
}

func (e *InvalidUserError) Error() string {
	return e.cause.Error()
}

var errInvalidPasswordLength = fmt.Errorf("password must be between %d and %d characters long", minPasswordLen, maxPasswordLen)
