package rest

import (
	"encoding/json"
	"github.com/angusgmorrison/typeddtalk/domain/user"
	"io"
	"net/http"
)

// userHandler handles all requests to /users.
type userHandler struct {
	service user.Service
}

type createUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// toDomain attempts to construct a valid [CreateUserRequest] from the request body, propagating any errors.
func (reqBody *createUserRequestBody) toDomain() (user.CreateUserRequest, error) {
	email, err := user.NewEmailAddress(reqBody.Email)
	if err != nil {
		return user.CreateUserRequest{}, err
	}

	passwordHash, err := user.NewPasswordHash(reqBody.Password)
	if err != nil {
		return user.CreateUserRequest{}, err
	}

	return user.NewCreateUserRequest(email, passwordHash, user.Bio(reqBody.Bio)), nil
}

type createUserResponseBody struct {
	ID string `json:"id"`
}

// POST /users
func (handler *userHandler) create(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a valid user representation.
	userReq, err := parseCreateUserRequest(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Create a user from valid inputs.
	user, err := handler.service.Create(userReq)
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

// parseCreateUserRequest attempts to parse a [user.CreateUserRequest] from the request body, propagating any errors.
func parseCreateUserRequest(r io.Reader) (user.CreateUserRequest, error) {
	var reqBody createUserRequestBody
	if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
		return user.CreateUserRequest{}, err
	}

	return reqBody.toDomain()
}

func handleError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *user.ParseError, *user.ConstraintViolationError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		panic(err) // 500
	}
}
