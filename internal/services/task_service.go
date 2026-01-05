package services

import (
	"context"
	"strconv"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskService struct {
	Repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		Repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {

	err := s.Repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTasks(ctx context.Context, filters map[string]string) ([]models.Task, error) {

	filter := bson.M{}

	if v, ok := filters["category"]; ok && v != "" {
		filter["category"] = v
	}
	if v, ok := filters["status"]; ok && v != "" {
		filter["status"] = v
	}

	if v, ok := filters["search"]; ok && v != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": v, "$options": "i"}},
			{"description": bson.M{"$regex": v, "$options": "i"}},
		}
	}

	sort := bson.D{}
	if v, ok := filters["sort"]; ok && v != "" {
		order := 1
		if v2, ok2 := filters["order"]; ok2 && v2 == "desc" {
			order = -1
		}
		sort = append(sort, bson.E{Key: v, Value: order})
	}

	limit := 20
	if v, ok := filters["limit"]; ok && v != "" {
		limit, _ = strconv.Atoi(v)
	}

	skip := 0
	if v, ok := filters["page"]; ok && v != "" {
		page, _ := strconv.Atoi(v)
		if page > 1 {
			skip = (page - 1) * limit
		}
	}

	tasks, err := s.Repo.GetTasks(ctx, filter, sort, limit, skip)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id primitive.ObjectID, req models.UpdateTaskRequest) (*models.Task, error) {

	task, err := s.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = *req.DueDate
	}

	updatedTask, err := s.Repo.UpdateTask(ctx, task)

	return updatedTask, err

}
