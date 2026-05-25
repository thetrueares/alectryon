package entities

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type HistoryModel struct {
	ID        bson.ObjectID `bson:"_id"`
	User      string        `bson:"user"`
	Direction string        `bson:"direction"`
	Message   string        `bson:"message"`
	Response  string        `bson:"response"`
	CreatedAt *time.Time    `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}

func NewHistoryRepository(collection *mongo.Collection) *HistoryRepository {
	return &HistoryRepository{collection: collection}
}

type HistoryRepository struct {
	collection *mongo.Collection
}

func (hr HistoryRepository) Save(history HistoryModel) error {

	now := time.Now()
	if history.CreatedAt == nil {
		history.CreatedAt = &now
	}

	history.UpdatedAt = now
	opts := options.UpdateOne().SetUpsert(true)
	_, err := hr.collection.UpdateOne(context.TODO(), bson.M{"_id": history.ID}, bson.D{{"$set", history}}, opts)

	if err != nil {
		return err
	}

	return nil
}

func NewInwardMessage(user, message string) HistoryModel {
	return newMessage(user, message, "inward")
}

func NewOutwardMessage(user, message string) HistoryModel {
	return newMessage(user, message, "outward")
}

func newMessage(user, message, direction string) HistoryModel {
	return HistoryModel{
		ID:        bson.NewObjectID(),
		User:      user,
		Direction: direction,
		Message:   message,
	}
}
