package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InputType string

const (
	InputTypeTelegramBot InputType = "telegram"
	InputTypeSlackBot    InputType = "slack"
	InputTypeAudio       InputType = "audio"
	InputTypeVideo       InputType = "video"
)

type InputModel struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	Type      InputType     `bson:"type"`
	Active    bool          `bson:"active"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}

type InputRepository struct {
	collection *mongo.Collection
}

func (ir InputRepository) GetAll() ([]InputModel, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cursor, err := ir.collection.Find(context.TODO(), nil, opts)

	if err != nil {
		return nil, err
	}

	var results []InputModel
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func NewInputRepository(collection *mongo.Collection) *InputRepository {
	return &InputRepository{collection: collection}
}
