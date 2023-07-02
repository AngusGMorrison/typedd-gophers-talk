package rest

import (
	"encoding/json"
	"github.com/angusgmorrison/typeddtalk/domain"
	"io"
	"net/http"
)

// userHandler handles all requests to /users.
type userHandler struct {
	service domain.UserService
}

type createUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// toDomain attempts to construct a valid [CreateUserRequest] from the request body, propagating any errors.
func (reqBody *createUserRequestBody) toDomain() (domain.CreateUserRequest, error) {
	email, err := domain.NewEmailAddress(reqBody.Email)
	if err != nil {
		return domain.CreateUserRequest{}, err
	}

	passwordHash, err := domain.NewPasswordHash(reqBody.Password)
	if err != nil {
		return domain.CreateUserRequest{}, err
	}

	return domain.NewCreateUserRequest(email, passwordHash, domain.Bio(reqBody.Bio)), nil
}

type createUserResponseBody struct {
	ID string `json:"id"`
}

// POST /users
func (uh *userHandler) create(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a valid domain representation.
	domainReq, err := parseCreateUserRequest(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Create a user from valid inputs.
	user, err := uh.service.Create(domainReq)
	if err != nil {
		handleError(w, err)
		return
	}

	// Respond with the user's ID.
	resBody := createUserResponseBody{ID: user.ID()}
	b, err := json.Marshal(&resBody)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}

// parseCreateUserRequest attempts to parse a [domain.CreateUserRequest] from the request body, propagating any errors.
func parseCreateUserRequest(r io.Reader) (domain.CreateUserRequest, error) {
	var reqBody createUserRequestBody
	if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
		return domain.CreateUserRequest{}, err
	}

	return reqBody.toDomain()
}

func handleError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *domain.ParseError, *domain.ConstraintViolationError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		panic(err) // 500
	}
}
