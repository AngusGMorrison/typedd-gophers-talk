package users

// Service is the interface to the business logic.
type Service interface {
	// Create and return a valid user, or an error if the user could not be created.
	Create(req CreateUserRequest) (User, error)
}

// Repository represents a store of user data.
type Repository interface {
	// Create persists a new user, returning the created [User] on success.
	// MUST return [ConstraintViolationError] if any field violates a repository constraint.
	Create(req CreateUserRequest) (User, error)
}
