package typedd

// Option represents a value that must be complete if it is present, but may be absent.
type Option[T Complete] struct {
	value T
	some  bool
}

// Some returns an Option[T] containing the given value if it is complete, or an error if it is incomplete.
func Some[T Complete](value T) (Option[T], error) {
	if !value.Complete() {
		return Option[T]{}, &IncompleteTypeError{Incomplete: value} // error at instantiation
	}

	return Option[T]{
		value: value,
		some:  true,
	}, nil
}

// None returns an empty, zero-valued for Option[T]. That is, the zero-value for Option is a valid value.
func None[T Complete]() Option[T] {
	return Option[T]{}
}

// Value returns the value of the Option[T] and true if it is Some, or T's zero-value and false if it is None.
func (o *Option[T]) Value() (T, bool) {
	return o.value, o.some
}
