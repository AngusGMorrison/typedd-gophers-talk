package typedd

import (
	"fmt"
	"reflect"
)

// ZeroAware types report whether they are allowed to be zero.
type ZeroAware interface {
	// Zeroable MUST return true if the zero-value for the type is a valid value.
	Zeroable() bool
}

// ZeroValueError indicates that a value which may never be zero is zero.
type ZeroValueError struct {
	ZeroValue ZeroAware
}

func (e *ZeroValueError) Error() string {
	return fmt.Sprintf("found zero-value of %T where non-zero value is required", e.ZeroValue)
}

// ValidateNonZero is a convenience function for validating that multiple NonZero values are indeed non-zero.
func ValidateNonZero(maybeZeroes ...ZeroAware) error {
	for _, maybeZero := range maybeZeroes {
		if !maybeZero.Zeroable() && reflect.ValueOf(maybeZero).IsZero() { // yuck
			return &ZeroValueError{ZeroValue: maybeZero}
		}
	}

	return nil
}
