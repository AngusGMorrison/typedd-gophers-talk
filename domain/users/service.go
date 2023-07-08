package users

// service is the implementation of Service.
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create a new user.
// Returns [ConstraintViolationError] if any field violates a repository constraint.
func (s *service) Create(req *CreateUserRequest) (*User, error) {
	user, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
