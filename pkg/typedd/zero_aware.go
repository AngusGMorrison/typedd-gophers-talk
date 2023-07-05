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

// MustBeNonZero is a convenience function for asserting that values which must be non-zero are indeed non-zero.
// Panics if a zero value for a non-zeroable type is encountered.
func MustBeNonZero(maybeZeroes ...ZeroAware) {
	for _, maybeZero := range maybeZeroes {
		mustBeNonZero(maybeZero)
	}
}

// Recursively asserts that no constituent ZeroAware field or element of maybeZero is zero.
func mustBeNonZero(maybeZero ZeroAware) {
	value := reflect.ValueOf(maybeZero)
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if zeroAware, ok := value.Field(i).Interface().(ZeroAware); ok {
				mustBeNonZero(zeroAware)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if zeroAware, ok := value.Index(i).Interface().(ZeroAware); ok {
				mustBeNonZero(zeroAware)
			}
		}
	case reflect.Ptr:
		if zeroAware, ok := value.Elem().Interface().(ZeroAware); ok {
			mustBeNonZero(zeroAware)
		}
	default:
		return // kind is not zero-aware; do nothing
	}

	if !maybeZero.Zeroable() && value.IsZero() {
		panic(fmt.Errorf("non-zeroable type %T had zero value", maybeZero))
	}
}
