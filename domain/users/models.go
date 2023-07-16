package users

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// EmailAddress represents an RFC 5322-compliant email address.
type EmailAddress struct {
	Raw string
}

// NewEmailAddress parses a Raw email address according to RFC 5322.
// Returns [ParseError] if the address is invalid.
func NewEmailAddress(raw string) (EmailAddress, error) {
	_, err := mail.ParseAddress(raw)
	if err != nil {
		return EmailAddress{}, NewParseRFC5233EmailError(err)
	}
	return EmailAddress{Raw: raw}, nil
}

func (e EmailAddress) String() string {
	return e.Raw
}

// PasswordHash represents a bcrypt-hashed password.
type PasswordHash struct {
	Bytes []byte
}

const (
	minPasswordLen = 8
	maxPasswordLen = 72 // bcrypt max length
)

// NewPasswordHash checks that rawPassword matches length constraints and hashes it with bcrypt.
// Returns [ParseError] if any stage fails.
func NewPasswordHash(rawPassword string) (PasswordHash, error) {
	if len(rawPassword) < minPasswordLen || len(rawPassword) > maxPasswordLen {
		return PasswordHash{}, NewPasswordLengthError()
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return PasswordHash{}, NewHashPasswordError(err)
	}

	return PasswordHash{Bytes: bytes}, nil
}

func (p PasswordHash) String() string {
	return string(p.Bytes)
}

// Bio represents a user's biography, which may be empty.
type Bio string

// User is a domain model representing a service user. It must be valid at all costs!
type User struct {
	ID           uuid.UUID
	Email        EmailAddress // required
	PasswordHash PasswordHash // required
	Bio          Bio          // optional
}

// NewUser creates a new, valid [User] with the given fields.
func NewUser(id uuid.UUID, email EmailAddress, passwordHash PasswordHash, bio Bio) *User {
	return &User{ID: id, Email: email, PasswordHash: passwordHash, Bio: bio}
}

// CreateUserRequest contains the valid fields required to create a new [User].
type CreateUserRequest struct {
	Email        EmailAddress
	PasswordHash PasswordHash
	Bio          Bio
}

func NewCreateUserRequest(email EmailAddress, passwordHash PasswordHash, bio Bio) *CreateUserRequest {
	return &CreateUserRequest{
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio,
	}
}
