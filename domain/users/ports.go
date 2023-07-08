package users

// Service is the interface to the business logic.
type Service interface {
	// Create and return a valid user, or an error if the user could not be created.
	Create(email, password, bio string) (*User, error)
}

// Repository represents a store of user data.
type Repository interface {
	// Create should persist a new user, returning the created [User] on success, or the User
	// zero-value and an error if the operation fails.
	// MUST return [ConstraintViolationError] if any field violates a repository constraint.
	Create(email, bio string, passwordHash []byte) (*User, error)
}
