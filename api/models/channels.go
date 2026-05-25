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

type ChannelEntity struct {
	ID        bson.ObjectID  `bson:"_id"`
	Name      string         `bson:"name"`
	Type      InputType      `bson:"type"`
	Active    bool           `bson:"active"`
	Options   map[string]any `bson:"options"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at,omitempty"`
}

func NewChannelRepository(collection *mongo.Collection) *ChannelRepository {
	return &ChannelRepository{collection: collection}
}

type ChannelRepository struct {
	collection *mongo.Collection
}

func (ir ChannelRepository) Save(input ChannelEntity) error {

	opts := options.UpdateOne().SetUpsert(true)
	_, err := ir.collection.UpdateOne(context.TODO(), bson.M{"_id": input.ID}, bson.D{{"$set", input}}, opts)

	if err != nil {
		return err
	}

	return nil
}

func (ir ChannelRepository) GetAll() ([]ChannelEntity, error) {
	cursor, err := ir.collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}

	var results []ChannelEntity
	if err = cursor.All(context.TODO(), &results); err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return results, nil
}

func (ir ChannelRepository) GetById(id string) (ChannelEntity, error) {
	var result ChannelEntity
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = ir.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)

	return result, err
}

func (ir ChannelRepository) DeleteById(id string) error {

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = ir.collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
