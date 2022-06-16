package config

import (
	"context"
	"log"

	"github.com/romik1505/ApiGateway/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoConnection(ctx context.Context, connString string) store.Storage {
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatalln("database connection err: %w", err)
	}

	if err := con.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln("database ping error: %w", err)
	}

	return store.Storage{
		Client: con,
	}
}
