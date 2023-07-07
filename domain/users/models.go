package users

import (
	"github.com/angusgmorrison/typeddtalk/pkg/typedd"
	"github.com/google/uuid"
	"net/mail"
)

// UUID represents a valid, universally unique identifier that cannot be zero.
type UUID struct {
	inner uuid.UUID
}

// NewUUID creates a new [UUID] from a raw string.
func NewUUID(raw string) (UUID, error) {
	inner, err := uuid.FromBytes([]byte(raw))
	if err != nil {
		return UUID{}, NewParseUUIDError(err)
	}

	return UUID{inner: inner}, nil
}

func (id UUID) Complete() bool {
	return id.inner != uuid.Nil
}

func (id UUID) String() string {
	return id.inner.String()
}

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

func (ea EmailAddress) Complete() bool {
	return ea.raw != ""
}

func (ea EmailAddress) String() string {
	return ea.raw
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

func (p PasswordHash) Complete() bool {
	return len(p.bytes) > 0
}

func (p PasswordHash) String() string {
	return string(p.bytes)
}

// Bio represents a user's biography, which may be empty.
type Bio string

// User is a domain model representing a service user. It must be valid at all costs!
type User struct {
	id           UUID
	email        EmailAddress // required
	passwordHash PasswordHash // required
	bio          Bio          // optional
}

// NewUser creates a new, valid [User] with the given fields.
func NewUser(id UUID, email EmailAddress, passwordHash PasswordHash, bio Bio) (User, error) {
	user := User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}

	if err := typedd.ValidateCompleteness(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *User) Complete() bool {
	return u.id.Complete() && u.email.Complete() && u.passwordHash.Complete()
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

func NewCreateUserRequest(email EmailAddress, passwordHash PasswordHash, bio Bio) (CreateUserRequest, error) {
	req := CreateUserRequest{
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}

	if err := typedd.ValidateCompleteness(&req); err != nil {
		return CreateUserRequest{}, err
	}

	return req, nil
}

func (req *CreateUserRequest) Complete() bool {
	return req.email.Complete() && req.passwordHash.Complete()
}
