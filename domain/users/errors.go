package users

import "fmt"

// Field represents a field of a domain type.
type Field int

const (
	EmailField Field = iota + 1
	PasswordHashField
	BioField
)

var fieldStrings = [3]string{"email", "password hash", "bio"}

func (f Field) String() string {
	return fieldStrings[f-1]
}

// ParseError indicates failure to parse a raw value into a valid domain type.
type ParseError struct {
	field    Field
	messages []string
	cause    error
}

func (e *ParseError) Field() Field {
	return e.field
}

func (e *ParseError) Messages() []string {
	return e.messages
}

func (e *ParseError) Cause() error {
	return e.cause
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid %s: %v", e.field, e.messages)
}

func NewParseRFC5233EmailError(cause error) *ParseError {
	return &ParseError{
		field:    EmailField,
		messages: []string{"email address is not RFC 5322-compliant"},
		cause:    cause,
	}
}

var passwordLengthErrorMessage = fmt.Sprintf("password must be between %d and %d characters", minPasswordLen, maxPasswordLen)

func NewPasswordLengthError() *ParseError {
	return &ParseError{
		field:    PasswordHashField,
		messages: []string{passwordLengthErrorMessage},
		cause:    nil,
	}
}

func NewHashPasswordError(cause error) *ParseError {
	return &ParseError{
		field:    PasswordHashField,
		messages: []string{"invalid password"},
		cause:    cause,
	}
}

// ConstraintViolationError indicates that a field has violated a repository constraint, such as a unique index.
type ConstraintViolationError struct {
	field    Field
	messages []string
	cause    error
}

func (e *ConstraintViolationError) Error() string {
	return fmt.Sprintf("constraint violation on %s: %v", fieldStrings[e.field], e.messages)
}
