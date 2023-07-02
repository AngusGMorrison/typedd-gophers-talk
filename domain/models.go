package domain

import (
	"github.com/google/uuid"
	"net/mail"
)

// EmailAddress represents an RFC 5322-compliant email address.
type EmailAddress struct {
	raw string
}

// NewEmailAddress parses a raw email address according to RFC 5322.
// Returns [ParseError] if the address is invalid.
func NewEmailAddress(raw string) (EmailAddress, error) {
	_, err := mail.ParseAddress(raw)
	if err != nil {
		return EmailAddress{}, NewParseRFC5233EmailError(err)
	}
	return EmailAddress{raw: raw}, nil
}

func (e EmailAddress) String() string {
	return e.raw
}

// PasswordHash represents a bcrypt-hashed password.
type PasswordHash struct {
	bytes []byte
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

	return PasswordHash{bytes: []byte(rawPassword)}, nil
}

func (p PasswordHash) String() string {
	return string(p.bytes)
}

// Bio represents a user's biography, which may be empty.
type Bio string

// User is a domain model representing a service user. It must be valid at all costs!
type User struct {
	id           uuid.UUID
	email        EmailAddress // required
	passwordHash PasswordHash // required
	bio          Bio          // optional
}

// NewUser creates a new, valid [User] with the given fields.
func NewUser(id uuid.UUID, email EmailAddress, passwordHash PasswordHash, bio Bio) User {
	return User{id: id, email: email, passwordHash: passwordHash, bio: bio}
}

func (u *User) ID() string {
	return u.id.String()
}

func (u *User) Email() EmailAddress {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) Bio() Bio {
	return u.bio
}

// CreateUserRequest contains the valid fields required to create a new [User].
type CreateUserRequest struct {
	email        EmailAddress
	passwordHash PasswordHash
	bio          Bio
}

func NewCreateUserRequest(email EmailAddress, passwordHash PasswordHash, bio Bio) CreateUserRequest {
	return CreateUserRequest{
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}
}
