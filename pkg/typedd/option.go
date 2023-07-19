package typedd

import "reflect"

// Option represents an optional value.
// The zero-value of T ALWAYS considered valid, and is typically used to indicate the removal of an existing value.
// Where T: [Complete], partially complete values are considered invalid.
type Option[T any] struct {
	value T
	some  bool
}

// None returns an empty, zero-valued for Option[T]. Hence, the zero-value of Option is always valid.
func None[T Complete]() Option[T] {
	return Option[T]{}
}

// Some returns an Option[T] containing the given value, which may be T's zero value.
// Where T: [Complete], invoking Some with a partially-complete value returns [IncompleteTypeError].
func Some[T any](value T) (Option[T], error) {
	if complete, ok := any(value).(Complete); ok {
		if partiallyComplete(complete) {
			return Option[T]{}, &IncompleteTypeError{Incomplete: complete}
		}
	}

	return Option[T]{
		value: value,
		some:  true,
	}, nil
}

// Value returns the value of the Option[T] and true if it is Some, or T's zero-value and false if it is None.
func (o *Option[T]) Value() (T, bool) {
	return o.value, o.some
}

func partiallyComplete(maybePartial Complete) bool {
	return !reflect.ValueOf(maybePartial).IsZero() && !maybePartial.Complete()
}
