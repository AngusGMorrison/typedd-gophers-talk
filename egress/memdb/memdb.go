package memdb

import (
	"errors"
	"fmt"
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"github.com/google/uuid"
)

type userRecord struct {
	id           users.UUID
	email        users.EmailAddress
	passwordHash users.PasswordHash
	bio          users.Bio
}

// ThreadUnsafeMemDB is a simple in-memory database implementation that is not thread-safe.
// Satisfies [users.Repository].
type ThreadUnsafeMemDB struct {
	usersByID    map[users.UUID]*userRecord
	usersByEmail map[users.EmailAddress]*userRecord
}

var _ users.Repository = (*ThreadUnsafeMemDB)(nil)

// NewThreadUnsafeMemDB creates a new [ThreadUnsafeMemDB].
func NewThreadUnsafeMemDB() *ThreadUnsafeMemDB {
	return &ThreadUnsafeMemDB{
		usersByID:    make(map[users.UUID]*userRecord),
		usersByEmail: make(map[users.EmailAddress]*userRecord),
	}
}

var errUserExists = errors.New("user already exists")

// Create a new user.
func (db *ThreadUnsafeMemDB) Create(req users.CreateUserRequest) (users.User, error) {
	if _, ok := db.usersByEmail[req.Email()]; ok {
		return users.User{}, &users.ConstraintViolationError{
			Field:    users.EmailField,
			Messages: []string{fmt.Sprintf("users with email %q already exists", req.Email())},
			Cause:    errUserExists,
		}
	}

	id, err := users.NewUUID(uuid.New())
	if err != nil {
		return users.User{}, err
	}

	record := userRecord{
		id:           id,
		email:        req.Email(),
		passwordHash: req.PasswordHash(),
		bio:          req.Bio(),
	}

	db.usersByID[record.id] = &record
	db.usersByEmail[record.email] = &record

	return users.NewUser(id, req.Email(), req.PasswordHash(), req.Bio()), nil
}

// Update the user identified by req.ID.
func (db *ThreadUnsafeMemDB) Update(req users.UpdateUserRequest) error {
	record, ok := db.usersByID[req.UserID()]
	if !ok {
		return &users.NotFoundError{UserID: req.UserID()}
	}

	// Return values from function calls like this getter invocation are not addressable, so we must save them to a
	// local variable before performing further work with them. We could return a pointer to the value from the getter,
	// but this would allow the caller to mutate the value.
	email := req.Email()
	if newEmail, ok := email.Value(); ok {
		delete(db.usersByEmail, record.email)
		record.email = newEmail
		db.usersByEmail[record.email] = record
	}

	passwordHash := req.PasswordHash()
	if newPasswordHash, ok := passwordHash.Value(); ok {
		record.passwordHash = newPasswordHash
	}

	bio := req.Bio()
	if newBio, ok := bio.Value(); ok {
		record.bio = newBio
	}

	return nil
}
