package users

// Service is the interface to the business logic.
type Service interface {
	// Create and return a valid users, or an error if the users could not be created.
	Create(email, password, bio string) (*User, error)
}

// Repository represents a store of users data.
type Repository interface {
	// Create persists a new users, returning the created [User] on success.
	Create(email, bio string, passwordHash []byte) (*User, error)
}
