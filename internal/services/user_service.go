package services

import (
	"context"
	"fmt"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/repository"
	"github.com/nv-root/task-manager/internal/utils"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.CreateUserRequest) (*models.UserResponse, error) {

	existingUser, _ := s.Repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return nil, utils.BadRequest("Email already exists", nil)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, utils.Internal("Error creating user", nil)
	}

	newUser := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	err = s.Repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, creds *models.Credentials) (any, error) {
	user, _ := s.Repo.GetUserByEmail(ctx, creds.Email)
	if user == nil {
		return nil, utils.Unauthorized("Invalid email or password", nil)
	}

	err := utils.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return nil, utils.Unauthorized("Invalid email or password", nil)
	}

	fmt.Printf("DEBUG: userId in service: %v\n", user.ID)

	token, err := utils.CreateTokenWithClaims(*user)
	if err != nil {
		fmt.Printf("DEBUG: Error logging in %v\n", err)
		return nil, utils.Internal("Error Logging In", nil)
	}

	return struct {
		Token string              `json:"token"`
		User  models.UserResponse `json:"user"`
	}{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
