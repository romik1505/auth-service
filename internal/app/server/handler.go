package server

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/romik1505/ApiGateway/internal/app/mapper"
	"github.com/romik1505/ApiGateway/internal/app/service"
)

type Handler struct {
	authService *service.AuthService
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := r.Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	req := mapper.LoginRequest{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := r.Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	req := mapper.TokenPair{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, service.ErrTokensNotFormPair) ||
			errors.Is(err, service.ErrTokensExpired) ||
			errors.Is(err, service.ErrInvalidToken) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
