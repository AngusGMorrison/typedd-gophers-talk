//go:build debug

package users

import "github.com/angusgmorrison/typeddtalk/pkg/typedd"

// nonZeroService implements [Service] by decorating the release [service] to assert that inputs types that are
// required to be non-zero are actually non-zero.
type nonZeroService struct {
	releaseService Service
}

// NewService decorates a [service] with non-zero assertions for [Service] input and [Repository] return types.
func NewService(repo Repository) Service {
	return &nonZeroService{
		releaseService: newService(
			&nonZeroRepository{releaseRepo: repo},
		),
	}
}

// Create panics if any non-zeroable input, as defined by its implementation of [typedd.ZeroAware] is zero.
func (s *nonZeroService) Create(req CreateUserRequest) (User, error) {
	typedd.MustBeNonZero(&req)
	return s.releaseService.Create(req)
}

// nonZeroRepository implements [Repository] by decorating the release [repository] to validate that outputs are
// non-zero.
type nonZeroRepository struct {
	releaseRepo Repository
}

// Create panics if any non-zeroable output, as defined by its implementation of [typedd.ZeroAware] is zero.
func (r *nonZeroRepository) Create(req CreateUserRequest) (User, error) {
	user, err := r.releaseRepo.Create(req)
	if err != nil {
		return User{}, err
	}

	typedd.MustBeNonZero(&user)

	return user, nil
}
