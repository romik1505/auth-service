package server

import (
	"github.com/gorilla/mux"
	_ "github.com/romik1505/ApiGateway/docs"
	"github.com/romik1505/ApiGateway/internal/app/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(s *service.AuthService) *mux.Router {
	router := mux.NewRouter()
	h := Handler{
		authService: s,
	}

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/login", h.login).Methods("POST")
	router.HandleFunc("/refresh-token", h.refreshToken).Methods("POST")

	return router
}
