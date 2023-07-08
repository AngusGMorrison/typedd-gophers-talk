package main

import (
	"github.com/angusgmorrison/typeddtalk/domain/users"
	"github.com/angusgmorrison/typeddtalk/egress/memdb"
	"github.com/angusgmorrison/typeddtalk/ingress/rest"
)

// Runs a crude HTTP server on 8080, exposing
//   - POST /users
//   - PUT /users
//
// Has no logging and no means of recovering from panics. Its in-memory DB is not thread-safe.
func main() {
	repo := memdb.NewThreadUnsafeMemDB()
	service := users.NewService(repo)
	server := rest.NewServer(":8080", service)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
