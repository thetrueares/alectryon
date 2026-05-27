package entities

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserEntity struct {
	ID           bson.ObjectID `bson:"_id"`
	Name         string        `bson:"name"`
	UserChannels []UserChannel `bson:"user_channels"`
	CreatedAt    *time.Time    `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

type UserChannel struct {
	ChannelType ChannelType   `bson:"channel_type"`
	ChannelID   bson.ObjectID `bson:"channel_id"`
	UserID      string        `bson:"user_id"`
}

type EmbeddedUser struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

type UserRepository struct {
	collection *mongo.Collection
}

func (ur UserRepository) Save(user UserEntity) error {

	now := time.Now()
	if user.CreatedAt == nil {
		user.CreatedAt = &now
	}

	user.UpdatedAt = now
	opts := options.UpdateOne().SetUpsert(true)
	_, err := ur.collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{{"$set", user}}, opts)

	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) FindByChannelSender(channelType ChannelType, userId string) (*UserEntity, error) {
	var user UserEntity
	filter := bson.M{
		"user_channels": bson.M{
			"$elemMatch": bson.M{
				"channel_type": channelType,
				"user_id":      userId,
			},
		},
	}

	err := ur.collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
		return nil, NoEntityFound{fmt.Sprintf("Can't find a user for channel %s and user id %s", channelType, userId)}
	}

	return &user, nil
}
