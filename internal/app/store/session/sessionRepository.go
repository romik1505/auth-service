package session

import (
	"context"
	"fmt"

	"github.com/romik1505/ApiGateway/internal/app/model"
	"github.com/romik1505/ApiGateway/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
)

type SessionRepository struct {
	Storage store.Storage
}

func NewSessionRepository(storage store.Storage) ISessionRepository {
	return &SessionRepository{
		Storage: storage,
	}
}

type ISessionRepository interface {
	CreateSession(ctx context.Context, session model.RefreshSession) error
	DeleteSession(ctx context.Context, userID string, sessionID string) (model.RefreshSession, error)
	GetSession(ctx context.Context, userID string, sessionID string) (model.RefreshSession, error)
}

func (s SessionRepository) CreateSession(ctx context.Context, session model.RefreshSession) error {
	_, err := s.Storage.Database("apigw").Collection("refresh-sessions").InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("creation session failed %v", err)
	}
	return nil
}

func (s SessionRepository) DeleteSession(ctx context.Context, userID string, sessionID string) (model.RefreshSession, error) {
	var session model.RefreshSession
	filter := bson.D{{"_id", sessionID}, {"user_id", userID}}
	res := s.Storage.Database("apigw").Collection("refresh-sessions").FindOneAndDelete(ctx, filter)

	return session, res.Decode(&session)
}

func (s SessionRepository) GetSession(ctx context.Context, userID string, sessionID string) (model.RefreshSession, error) {
	var result model.RefreshSession
	filter := bson.D{{"_id", sessionID}, {"user_id", userID}}

	err := s.Storage.Database("apigw").
		Collection("refresh-sessions").
		FindOne(ctx, filter).
		Decode(&result)
	return result, err
}
