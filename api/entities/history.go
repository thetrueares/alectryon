package entities

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type HistoryEntity struct {
	ID        bson.ObjectID `bson:"_id"`
	User      EmbeddedUser  `bson:"user"`
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

func (hr HistoryRepository) Save(history HistoryEntity) error {

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

func (hr HistoryRepository) GetLastFive() ([]HistoryEntity, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(5)
	cursor, err := hr.collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []HistoryEntity
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (hr HistoryRepository) GetLastTenForUser(user *UserEntity) ([]HistoryEntity, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(10)
	filter := bson.M{"user._id": user.ID}
	cursor, err := hr.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var results []HistoryEntity
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func NewInwardMessage(user *UserEntity, message string) HistoryEntity {
	return newMessage(user, message, "inward")
}

func NewOutwardMessage(user *UserEntity, message string) HistoryEntity {
	return newMessage(user, message, "outward")
}

func newMessage(user *UserEntity, message, direction string) HistoryEntity {
	embeddedUser := EmbeddedUser{
		ID:   user.ID,
		Name: user.Name,
	}

	return HistoryEntity{
		ID:        bson.NewObjectID(),
		User:      embeddedUser,
		Direction: direction,
		Message:   message,
	}
}
