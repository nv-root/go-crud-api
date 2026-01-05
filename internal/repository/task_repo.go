package repository

import (
	"context"
	"fmt"
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

func (tr *TaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	result := tr.Collection.FindOne(ctx, bson.M{"_id": id})

	var task models.Task
	err := result.Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (tr *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	filter := bson.M{"_id": task.ID}

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"due_date":    task.DueDate,
			"updated_at":  time.Now(),
		},
	}

	_, err := tr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// func (tr *TaskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {}

// func (tr *TaskRepository) MarkCompleted(ctx context.Context, id primitive.ObjectID) error {}
