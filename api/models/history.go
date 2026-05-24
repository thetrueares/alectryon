package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type History struct {
	ID        bson.ObjectID `bson:"_id"`
	Message   string        `bson:"message"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}
