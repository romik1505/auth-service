package store

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	*mongo.Client
}
