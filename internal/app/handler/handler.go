package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/romik1505/ApiGateway/docs"
	"github.com/romik1505/ApiGateway/internal/app/mapper"
	"github.com/romik1505/ApiGateway/internal/app/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/login", StartSpanMiddleware(ErrorWrapper(h.Login))).Methods("POST")
	router.HandleFunc("/refresh-token", StartSpanMiddleware(ErrorWrapper(h.RefreshToken))).Methods("POST")

	return router
}

// @Summary Login
// @Tags auth
// @Description login to service
// @ID login
// @Accept json
// @Produce json
// @Param input body mapper.LoginRequest true "account info"
// @Success 200 {object} mapper.TokenPair "token pairs"
// @Failure 400,500 {string} string
// @Router /login [post]
func (h Handler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := mapper.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return ErrorBadRequest
	}

	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return err
	}
	return nil
}

// @Summary Refresh token
// @Tags refresh
// @Description refresh token pair
// @ID refresh
// @Accept json
// @Produce json
// @Param input body mapper.TokenPair true "previous token pair"
// @Success 200 {object} mapper.TokenPair "next token pairs"
// @Failure 400,403,500 {string} string
// @Router /refresh-token [post]
func (h Handler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := mapper.TokenPair{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return ErrorBadRequest
	}

	resp, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrorBadRequest = errors.New("bad request")
)
