package rest

import (
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"net/http"
)

// Server is a crude HTTP server without logging or panic recovery.
type Server struct {
	*http.Server
}

func NewServer(addr string, userService users.Service) *Server {
	router := http.NewServeMux()
	router.Handle("/users", &usersHandler{service: userService})
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}
