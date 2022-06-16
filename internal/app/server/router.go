package server

import (
	"github.com/gorilla/mux"
	"github.com/romik1505/ApiGateway/internal/app/service"
)

func NewRouter(s *service.AuthService) *mux.Router {
	router := mux.NewRouter()
	h := Handler{
		authService: s,
	}
	router.HandleFunc("/login", h.login).Methods("POST")
	router.HandleFunc("/refresh-token", h.refreshToken).Methods("POST")

	return router
}
