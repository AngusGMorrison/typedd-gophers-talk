package users

import "github.com/angusgmorrison/typeddtalk/pkg/typedd"

// service is the implementation of Service.
type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

// Create a new users.
// Returns [ConstraintViolationError] if any field violates a repository constraint.
func (s *service) Create(req CreateUserRequest) (User, error) {
	if err := typedd.ValidateNonZero(&req); err != nil {
		return User{}, err
	}

	user, err := s.repo.Create(req)
	if err != nil {
		return User{}, err
	}
	if err = typedd.ValidateNonZero(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
