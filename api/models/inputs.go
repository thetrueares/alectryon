package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (ir InputRepository) Save(input InputModel) error {

	_, err := ir.collection.InsertOne(context.TODO(), input)

	if err != nil {
		return err
	}

	return nil
}

func (ir InputRepository) GetAll() ([]InputModel, error) {
	cursor, err := ir.collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}

	var results []InputModel
	if err = cursor.All(context.TODO(), &results); err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return results, nil
}

func NewInputRepository(collection *mongo.Collection) *InputRepository {
	return &InputRepository{collection: collection}
}
