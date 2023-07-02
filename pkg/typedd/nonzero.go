package typedd

import (
	"fmt"
	"reflect"
)

// ZeroValueError indicates that a value which may never be zero is zero.
type ZeroValueError struct {
	ZeroValue any
}

func (e *ZeroValueError) Error() string {
	return fmt.Sprintf("found zero-value of %T where non-zero value is required", e.ZeroValue)
}

// ValidateNonZero is a convenience function for validating that multiple NonZero values are indeed non-zero.
func ValidateNonZero(maybeZeroes ...any) error {
	for _, maybeZero := range maybeZeroes {
		if reflect.ValueOf(maybeZero).IsZero() { // yuck
			return &ZeroValueError{ZeroValue: maybeZero}
		}
	}

	return nil
}
