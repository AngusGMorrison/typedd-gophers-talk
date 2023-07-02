package rest

import (
	"encoding/json"
	"errors"
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"net/http"
)

// userHandler handles all requests to /users.
type userHandler struct {
	service users.Service
}

type createUserRequestBody struct {
	Email    string
	Password string
	Bio      string
}

type createUserResponseBody struct {
	ID string `json:"id"`
}

func (uh *userHandler) create(w http.ResponseWriter, r *http.Request) {
	var reqBody createUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	user, err := uh.service.Create(reqBody.Email, reqBody.Password, reqBody.Bio)
	if err != nil {
		if errors.Is(err, &users.InvalidUserError{}) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		panic(err) // let the recovery middleware handle this
	}

	resBody := createUserResponseBody{ID: user.ID.String()}
	b, err := json.Marshal(&resBody)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}
