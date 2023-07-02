package rest

import (
	"encoding/json"
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"io"
	"net/http"
)

// usersHandler handles all requests to /users.
type usersHandler struct {
	service users.Service
}

type createUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// toDomain attempts to construct a valid [CreateUserRequest] from the request body, propagating any errors.
func (reqBody *createUserRequestBody) toDomain() (users.CreateUserRequest, error) {
	email, err := users.NewEmailAddress(reqBody.Email)
	if err != nil {
		return users.CreateUserRequest{}, err
	}

	passwordHash, err := users.NewPasswordHash(reqBody.Password)
	if err != nil {
		return users.CreateUserRequest{}, err
	}

	req, err := users.NewCreateUserRequest(email, passwordHash, users.Bio(reqBody.Bio))
	if err != nil {
		return users.CreateUserRequest{}, err
	}

	return req, nil
}

type createUserResponseBody struct {
	ID string `json:"id"`
}

// POST /users
func (handler *usersHandler) create(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a valid User.
	userReq, err := parseCreateUserRequest(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Create a users from valid inputs.
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

// parseCreateUserRequest attempts to parse a [users.CreateUserRequest] from the request body, propagating any errors.
func parseCreateUserRequest(r io.Reader) (users.CreateUserRequest, error) {
	var reqBody createUserRequestBody
	if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
		return users.CreateUserRequest{}, err
	}

	return reqBody.toDomain()
}

func handleError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *users.ParseError, *users.ConstraintViolationError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		panic(err) // 500
	}
}
