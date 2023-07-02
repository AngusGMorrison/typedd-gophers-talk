package domain

// userService is the implementation of UserService.
type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// Create a new user.
// Returns [ConstraintViolationError] if any field violates a repository constraint.
func (s *userService) Create(req CreateUserRequest) (User, error) {
	user, err := s.repo.Create(req)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
