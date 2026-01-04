package repository

import (
	"context"
	"time"

	"github.com/nv-root/task-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TaskRepository struct {
	Collection *mongo.Collection
}

func NewTaskRespository(client *mongo.Client, dbName string) *TaskRepository {
	return &TaskRepository{
		Collection: client.Database(dbName).Collection("tasks"),
	}
}

func (tr *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err := tr.Collection.InsertOne(ctx, task)
	return err
}

func (tr *TaskRepository) GetTasks(ctx context.Context, filter bson.M, sort bson.D, limit, skip int) ([]models.Task, error) {
	opts := options.Find().SetSort(sort).SetLimit(int64(limit)).SetSkip(int64(skip))
	cursor, err := tr.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := []models.Task{}
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// func (tr *TaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Task, error) {}

// func (tr *TaskRepository) UpdateTask (ctx context.Context, t *Tasl) error {}

// func (tr *TaskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {}

// func (tr *TaskRepository) MarkCompleted(ctx context.Context, id primitive.ObjectID) error {}
