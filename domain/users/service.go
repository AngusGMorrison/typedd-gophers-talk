package users

import (
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// service is the implementation of Service.
type service struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &service{repo: repo}
}

// Create validates its inputs and creates a new [User], returning the saved user on success, or an error otherwise.
func (s *service) Create(email, password, bio string) (*User, error) {
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

const (
	minPasswordLen = 8
	maxPasswordLen = 74
)

func validatePassword(password string) error {
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return errInvalidPasswordLength
	}

	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
