package services

import (
	"context"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/repository"
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
