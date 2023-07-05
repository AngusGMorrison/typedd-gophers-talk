//go:build release

package users

// NewService returns a configured user [Service].
func NewService(repo Repository) {
	return newService(repo)
}
