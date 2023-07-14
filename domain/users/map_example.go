package users

type UserMetadata struct {
	m map[string]string
}

// NewVulnerableUserMetadata wraps the caller-provided map, trusting that they won't do anything nasty with their
// reference to it. Yuck!
func NewVulnerableUserMetadata(m map[string]string) UserMetadata {
	return UserMetadata{m: m}
}

// Pair is a simple key-value pair.
type Pair[T comparable, U any] struct {
	key T
	val U
}

func NewPair[T comparable, U any](key T, val U) Pair[T, U] {
	return Pair[T, U]{key: key, val: val}
}

// NewSafeUserMetadata requires the caller to deconstruct the original map into a slice of [Pair], passing this slice
// to our constructor as variadic arguments. We then transform this slice into a map under our control. Three full
// copies. PURITY ACHIEVED!
func NewSafeUserMetadata(m ...Pair[string, string]) UserMetadata {
	d := UserMetadata{
		m: make(map[string]string, len(m)),
	}
	for _, pair := range m {
		d.m[pair.key] = pair.val
	}
	return d
}
