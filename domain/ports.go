package domain

// UserService is the interface to the business logic.
type UserService interface {
	// Create and return a valid user, or an error if the user could not be created.
	Create(req CreateUserRequest) (User, error)
}

// UserRepository represents a store of user data.
// MUST return [ConstraintViolationError] if any field violates a repository constraint.
type UserRepository interface {
	Create(req CreateUserRequest) (User, error)
}
