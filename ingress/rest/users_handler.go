package rest

import (
	"encoding/json"
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"github.com/angusgmorrison/typeddtalk/pkg/typedd"
	"io"
	"net/http"
)

// usersHandler handles all requests to /users.
type usersHandler struct {
	service users.Service
}

// ServeHTTP satisfies [http.Handler].
func (handler *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.create(w, r)
	case http.MethodPatch:
		handler.update(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
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

	return users.NewCreateUserRequest(email, passwordHash, users.Bio(reqBody.Bio)), nil
}

type createUserResponseBody struct {
	ID string `json:"id"`
}

// POST /users
func (handler *usersHandler) create(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a valid users.CreateUserRequest.
	domainReq, err := parseCreateUserRequest(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Create a user from valid inputs.
	user, err := handler.service.Create(domainReq)
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

type updateUserRequestBody struct {
	ID       string  `json:"id"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Bio      *string `json:"bio"`
}

// PATCH /users
func (handler *usersHandler) update(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a valid users.UpdateUserRequest.
	domainReq, err := parseUpdateUserRequest(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Update the user with valid inputs.
	if err := handler.service.Update(domainReq); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// parseCreateUserRequest attempts to parse a [users.CreateUserRequest] from the request body, propagating any errors.
func parseCreateUserRequest(r io.Reader) (users.CreateUserRequest, error) {
	var reqBody createUserRequestBody
	if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
		return users.CreateUserRequest{}, err
	}

	return reqBody.toDomain()
}

func parseUpdateUserRequest(r io.Reader) (users.UpdateUserRequest, error) {
	var reqBody updateUserRequestBody
	if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
		return users.UpdateUserRequest{}, err
	}

	id, err := users.NewUUIDFromString(reqBody.ID)
	if err != nil {
		return users.UpdateUserRequest{}, err
	}

	var emailOpt typedd.Option[users.EmailAddress]
	if reqBody.Email != nil {
		email, err := users.NewEmailAddress(*reqBody.Email)
		if err != nil {
			return users.UpdateUserRequest{}, err
		}

		emailOpt, err = typedd.Some(email)
		if err != nil {
			return users.UpdateUserRequest{}, err
		}
	}

	var passwordOpt typedd.Option[users.PasswordHash]
	if reqBody.Password != nil {
		password, err := users.NewPasswordHash(*reqBody.Password)
		if err != nil {
			return users.UpdateUserRequest{}, err
		}

		passwordOpt, err = typedd.Some(password)
		if err != nil {
			return users.UpdateUserRequest{}, err
		}
	}

	var bioOpt typedd.Option[users.Bio]
	if reqBody.Bio != nil {
		bioOpt, err = typedd.Some(users.Bio(*reqBody.Bio))
		if err != nil {
			return users.UpdateUserRequest{}, err
		}
	}

	return users.NewUpdateUserRequest(id, emailOpt, passwordOpt, bioOpt), nil
}

func handleError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *users.ParseError, *users.ConstraintViolationError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		panic(err) // 500
	}
}
