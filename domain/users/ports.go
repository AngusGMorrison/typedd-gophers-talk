package users

// Service is the interface to the business logic.
type Service interface {
	// Create and return a valid user, or an error if the user could not be created.
	Create(req CreateUserRequest) (User, error)

	// Update the user identified by req.ID with the fields provided in req, returning an error if the requested update
	// cannot be performed.
	Update(req UpdateUserRequest) error
}

// Repository represents a store of user data.
type Repository interface {
	// Create should persist a new user, returning the created [User] on success, or return an error and the User
	// zero-value if the operation fails.
	// MUST return [ConstraintViolationError] if any field violates a repository constraint.
	Create(req CreateUserRequest) (User, error)

	// Update should update the fields of the user identified by req.ID for which a value is provided, or return an
	// error if the requested update cannot be performed.
	Update(req UpdateUserRequest) error
}
