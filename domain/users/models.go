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

// NewUUIDFromString creates a new [UUID] from a raw string.
func NewUUIDFromString(raw string) (UUID, error) {
	inner, err := uuid.Parse(raw)
	if err != nil {
		return UUID{}, NewParseUUIDError(err)
	}

	return UUID{inner: inner}, nil
}

// NewUUID creates a new [UUID] from a [uuid.UUID].
func NewUUID(id uuid.UUID) (UUID, error) {
	if id == uuid.Nil {
		return UUID{}, NewParseUUIDError(nil)
	}

	return UUID{inner: id}, nil
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

func (ea EmailAddress) Complete() bool {
	return ea.raw != ""
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

func (ph PasswordHash) Complete() bool {
	return len(ph.bytes) > 0
}

func (ph PasswordHash) String() string {
	return string(ph.bytes)
}

// Bio represents a user's biography, which may be empty.
type Bio string

func (b Bio) Complete() bool {
	return true
}

// User is a domain model representing a user account. It must be valid at all costs!
type User struct {
	id           UUID
	email        EmailAddress // required
	passwordHash PasswordHash // required
	bio          Bio          // optional
}

// NewUser creates a new, valid [User] with the given fields.
func NewUser(id UUID, email EmailAddress, passwordHash PasswordHash, bio Bio) User {
	return User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}
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

func NewCreateUserRequest(email EmailAddress, passwordHash PasswordHash, bio Bio) CreateUserRequest {
	return CreateUserRequest{
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}
}

func (req *CreateUserRequest) Complete() bool {
	return req.email.Complete() && req.passwordHash.Complete()
}

func (req *CreateUserRequest) Email() EmailAddress {
	return req.email
}

func (req *CreateUserRequest) PasswordHash() PasswordHash {
	return req.passwordHash
}

func (req *CreateUserRequest) Bio() Bio {
	return req.bio
}

// UpdateUserRequest contains the valid fields required to update an existing [User]. [typedd.Option] fields are used
// to distinguish between values that are empty because they should not be updated, and those which should be updated to
// an empty value.
type UpdateUserRequest struct {
	userID       UUID
	email        typedd.Option[EmailAddress]
	passwordHash typedd.Option[PasswordHash]
	bio          typedd.Option[Bio]
}

func NewUpdateUserRequest(
	id UUID,
	email typedd.Option[EmailAddress],
	passwordHash typedd.Option[PasswordHash],
	bio typedd.Option[Bio],
) UpdateUserRequest {
	return UpdateUserRequest{
		userID:       id,
		email:        email,
		passwordHash: passwordHash,
		bio:          bio,
	}
}

func (req *UpdateUserRequest) Complete() bool {
	return req.userID.Complete()
}

func (req *UpdateUserRequest) UserID() UUID {
	return req.userID
}

func (req *UpdateUserRequest) Email() typedd.Option[EmailAddress] {
	return req.email
}

func (req *UpdateUserRequest) PasswordHash() typedd.Option[PasswordHash] {
	return req.passwordHash
}

func (req *UpdateUserRequest) Bio() typedd.Option[Bio] {
	return req.bio
}
