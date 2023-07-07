//go:build debug

package users

import (
	"github.com/angusgmorrison/typeddtalk/pkg/typedd"
	"github.com/google/uuid"
)

// completeAssertingService implements [Service] by decorating the release [service] to assert that inputs types that are
// required to be non-zero are actually non-zero.
type completeAssertingService struct {
	releaseService Service
}

// NewService decorates a service with completeness assertions for [Service] input and [Repository] return types.
func NewService(repo Repository) Service {
	return &completeAssertingService{
		releaseService: newService(
			&completeAssertingRepository{releaseRepo: repo},
		),
	}
}

// Create panics if any input implementing [typedd.Complete] is incomplete.
func (s *completeAssertingService) Create(req CreateUserRequest) (User, error) {
	typedd.MustBeComplete(&req)
	return s.releaseService.Create(req)
}

// completeAssertingRepository implements [Repository] by decorating the release Repository to validate that outputs are
// non-zero.
type completeAssertingRepository struct {
	releaseRepo Repository
}

// Create panics if the [User] returned by the underlying [Repository.Create] implementation is incomplete.
func (r *completeAssertingRepository) Create(req CreateUserRequest) (User, error) {
	user, err := r.releaseRepo.Create(req)
	if err != nil {
		return User{}, err
	}

	typedd.MustBeComplete(&user)

	return user, nil
}

// Complete implementations.
func (ea EmailAddress) Complete() bool {
	return ea.raw != ""
}

func (ph PasswordHash) Complete() bool {
	return len(ph.bytes) > 0
}

func (id UUID) Complete() bool {
	return id.inner != uuid.Nil
}

func (u *User) Complete() bool {
	return u.id.Complete() && u.email.Complete() && u.passwordHash.Complete()
}

func (req *CreateUserRequest) Complete() bool {
	return req.email.Complete() && req.passwordHash.Complete()
}
