package users

import "github.com/google/uuid"

// User is a domain model representing a service users. It must be valid at all costs!
type User struct {
	ID           uuid.UUID
	Email        string // required
	PasswordHash []byte // required
	Bio          string // optional
}
