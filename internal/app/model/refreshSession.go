package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/romik1505/ApiGateway/internal/app/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	RefreshPrivateKey = []byte(config.GetValue(config.RefreshPrivateKey))
)

type RefreshSession struct {
	ID           string    `bson:"_id"`
	UserID       string    `bson:"user_id"`
	RefreshToken []byte    `bson:"refresh_token"`
	ExpiresIn    int64     `bson:"expires_in"`
	CreatedAt    time.Time `bson:"created_at"`
}

type RefreshTokenClaims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (r *RefreshSession) GenerateRefreshTokenString() (string, error) {
	claims := RefreshTokenClaims{
		ID:     r.ID,
		UserID: r.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: r.ExpiresIn,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(RefreshPrivateKey)
	if err != nil {
		return "", err
	}

	r.RefreshToken = []byte(signedToken)

	return signedToken, nil
}

func (r *RefreshSession) HashToken() error {
	hashedToken, err := bcrypt.GenerateFromPassword(r.RefreshToken, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.RefreshToken = hashedToken
	return nil
}
