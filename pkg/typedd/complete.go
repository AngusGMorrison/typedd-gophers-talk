package typedd

import (
	"fmt"
)

// Complete types require that all of their Complete fields are non-zero.
type Complete interface {
	// Complete MUST return false if:
	// * for primitive types, the value is the zero-value for the type; OR
	// * for complex types, calling Complete on any of its Complete fields return false.
	Complete() bool
}

// IncompleteTypeError indicates that a type satisfying Complete was incomplete.
type IncompleteTypeError struct {
	Incomplete Complete
}

func (e *IncompleteTypeError) Error() string {
	return fmt.Sprintf("value of type %[1]T implements Complete but was incomplete: %#[1]v", e.Incomplete)
}

// MustBeComplete panics if any of the given Complete types are incomplete.
func MustBeComplete(maybeComplete ...Complete) {
	for _, mc := range maybeComplete {
		if !mc.Complete() {
			panic(&IncompleteTypeError{Incomplete: mc})
		}
	}
}
