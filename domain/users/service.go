package users

// service is the implementation of Service.
type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

// Create a new user.
// Returns [ConstraintViolationError] if any field violates a repository constraint.
func (s *service) Create(req CreateUserRequest) (User, error) {
	return s.repo.Create(req)
}

// Update the user identified by req.ID.
func (s *service) Update(req UpdateUserRequest) error {
	return s.repo.Update(req)
}

// BulkUpdate the users contained in req.
func (s *service) BulkUpdate(req BulkUpdateUserRequest) error {
	return s.repo.BulkUpdate(req)
}
