package service

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/romik1505/ApiGateway/internal/app/mapper"
	"github.com/romik1505/ApiGateway/internal/app/model"
	"github.com/romik1505/ApiGateway/internal/app/store"
	"github.com/romik1505/ApiGateway/internal/app/store/session"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SessionRepository session.ISessionRepository
}

type IAuthService interface {
	Login(ctx context.Context, req mapper.LoginRequest) (mapper.TokenPair, error)
	RefreshToken(ctx context.Context, req mapper.TokenPair) (mapper.TokenPair, error)
}

func NewUserService(ctx context.Context, store store.Storage) *AuthService {
	return &AuthService{
		SessionRepository: session.NewSessionRepository(store),
	}
}

const (
	AccessTokenTTL  = time.Minute * 30
	RefreshTokenTTL = time.Hour * 24 * 30 // 30 days
)

var (
	errTokenNotValid     = errors.New("token not valid")
	errTokensNotFormPair = errors.New("tokens not form a pair")
	errTokensExpired     = errors.New("tokens expired")
)

func (a *AuthService) generateTokenPair(ctx context.Context, userID string) (mapper.TokenPair, error) {
	refreshSession := model.RefreshSession{
		ID:        uuid.NewString(),
		UserID:    userID,
		ExpiresIn: time.Now().Add(RefreshTokenTTL).Unix(),
		CreatedAt: time.Now(),
	}

	accessToken, err := model.NewSignedAccessToken(ctx, userID, refreshSession.ID)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	refreshToken, err := refreshSession.GenerateRefreshTokenString()
	if err != nil {
		return mapper.TokenPair{}, err
	}

	err = refreshSession.HashToken()
	if err != nil {
		return mapper.TokenPair{}, err
	}

	err = a.SessionRepository.CreateSession(ctx, refreshSession)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	return mapper.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: base64.StdEncoding.EncodeToString([]byte(refreshToken)),
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req mapper.LoginRequest) (mapper.TokenPair, error) {
	if err := req.Bind(); err != nil {
		return mapper.TokenPair{}, err
	}
	pair, err := a.generateTokenPair(ctx, req.UserID)
	if err != nil {
		return mapper.TokenPair{}, err
	}
	return pair, err
}

func (a *AuthService) RefreshToken(ctx context.Context, req mapper.TokenPair) (mapper.TokenPair, error) {
	if err := req.Bind(); err != nil {
		return mapper.TokenPair{}, err
	}

	token, err := jwt.ParseWithClaims(req.AccessToken, &model.JwtClaims{}, model.JWTKeyFunc)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	claims, ok := token.Claims.(*model.JwtClaims)
	if !ok || !token.Valid {
		return mapper.TokenPair{}, errTokenNotValid
	}

	refresh, err := base64.StdEncoding.DecodeString(req.RefreshToken)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	oldSession, err := a.SessionRepository.GetSession(ctx, claims.UserID, claims.RefreshTokenID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return mapper.TokenPair{}, errTokensNotFormPair
		}
		return mapper.TokenPair{}, err
	}

	if oldSession.ExpiresIn < time.Now().Unix() {
		return mapper.TokenPair{}, errTokensExpired
	}

	if err := bcrypt.CompareHashAndPassword(oldSession.RefreshToken, []byte(refresh)); err != nil {
		return mapper.TokenPair{}, errTokensNotFormPair
	}

	_, err = a.SessionRepository.DeleteSession(ctx, claims.UserID, claims.RefreshTokenID)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	res, err := a.generateTokenPair(ctx, claims.UserID)

	if err != nil {
		return mapper.TokenPair{}, err
	}
	return res, err
}
