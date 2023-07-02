package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// User is a domain model representing a service user. It must be valid at all costs!
type User struct {
	ID           uuid.UUID
	Email        string // required
	PasswordHash []byte // required
	Bio          string // optional
}

// UserService is the interface to the business logic.
type UserService interface {
	// Create and return a valid user, or an error if the user could not be created.
	Create(email, password, bio string) (*User, error)
}

// UserRepository represents a store of user data.
type UserRepository interface {
	Create(email, bio string, passwordHash []byte) (*User, error)
}

// userService is the implementation of UserService.
type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

type InvalidUserError struct {
	cause error
}

func (e *InvalidUserError) Error() string {
	return e.cause.Error()
}

func (s *userService) Create(email, password, bio string) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, &InvalidUserError{cause: err}
	}

	if err := validatePassword(password); err != nil {
		return nil, &InvalidUserError{cause: err}
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, &InvalidUserError{cause: err}
	}

	user, err := s.repo.Create(email, bio, passwordHash)
	if err != nil {
		return nil, &InvalidUserError{cause: err}
	}

	return user, nil
}

func validatePassword(password string) error {
	// enforce password rules
	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
