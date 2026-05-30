package entities

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TaskType string
type TaskStatus string

const (
	ProvideInformationTask TaskType   = "PROVIDE_INFORMATION"
	PerformActionTask      TaskType   = "PERFORM_ACTION"
	StatusNotStarted       TaskStatus = "STATUS_NOT_STARTED"
	StatusStarted          TaskStatus = "STATUS_STARTED"
	StatusCompleted        TaskStatus = "STATUS_COMPLETED"
	StatusFailed           TaskStatus = "STATUS_FAILED"
	StatusCancelled        TaskStatus = "STATUS_CANCELLED"
	StatusRestarted        TaskStatus = "STATUS_RESTARTED"
)

type TaskEntity struct {
	ID                  bson.ObjectID                              `bson:"_id,omitempty"`
	User                EmbeddedUser                               `bson:"user"`
	Status              string                                     `bson:"status"`
	Type                TaskType                                   `bson:"type"`
	RequiredInformation map[string]EmbeddedRequiredInformationData `bson:"required_information"`
	Description         string                                     `bson:"description"`
	TaskWorkOutput      []EmbeddedTaskWorkOutput                   `bson:"task_work_output"`
	CreatedAt           time.Time                                  `bson:"created_at"`
	UpdatedAt           time.Time                                  `bson:"updated_at"`
}
type EmbeddedTaskWorkOutput struct {
	WorkDone string `bson:"work_done"`
	Complete bool   `bson:"complete"`
	NextStep string `bson:"next_step"`
}
type EmbeddedTask struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Description string        `bson:"description"`
}
type EmbeddedRequiredInformationData struct {
	Required bool `bson:"required"`
	Value    any  `bson:"value"`
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{collection: collection}
}

type TaskRepository struct {
	collection *mongo.Collection
}

func (tr TaskRepository) Save(task *TaskEntity) error {

	now := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}

	task.UpdatedAt = now
	opts := options.UpdateOne().SetUpsert(true)
	_, err := tr.collection.UpdateOne(context.TODO(), bson.M{"_id": task.ID}, bson.D{{"$set", task}}, opts)

	if err != nil {
		return err
	}

	return nil
}

func (tr TaskRepository) FindById(id string) (*TaskEntity, error) {
	objectID, err := bson.ObjectIDFromHex(id)

	log.Printf("objectID: %v", objectID)
	if err != nil {
		return nil, err
	}

	var task TaskEntity
	err = tr.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
