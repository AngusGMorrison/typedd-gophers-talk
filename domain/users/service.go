package users

import "github.com/angusgmorrison/typeddtalk/pkg/typedd"

// service is the implementation of Service.
type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

// Create a new user.
// Returns [typedd.IncompleteTypeError] if req or the [User] returned by the repository violate completeness
// constraints.
// Returns [ConstraintViolationError] if any field violates a repository constraint.
func (s *service) Create(req CreateUserRequest) (User, error) {
	if err := typedd.ValidateCompleteness(&req); err != nil {
		return User{}, err
	}

	user, err := s.repo.Create(req)
	if err != nil {
		return User{}, err
	}
	if err = typedd.ValidateCompleteness(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
