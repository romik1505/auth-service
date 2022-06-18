package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/romik1505/ApiGateway/internal/app/mapper"
	"github.com/romik1505/ApiGateway/internal/app/service"
)

type Handler struct {
	authService *service.AuthService
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
func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := mapper.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
func (h Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := mapper.TokenPair{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, service.ErrTokensNotFormPair) ||
			errors.Is(err, service.ErrTokensExpired) ||
			errors.Is(err, service.ErrInvalidToken) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
