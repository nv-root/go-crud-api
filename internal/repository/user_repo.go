package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRespository(client *mongo.Client, dbName string) *UserRepository {
	return &UserRepository{
		Collection: client.Database(dbName).Collection("users"),
	}
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := ur.Collection.FindOne(ctx, bson.M{"email": email})
	err := result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		fmt.Printf("DEBUG: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.Collection.InsertOne(ctx, user)
	if err != nil {
		return utils.Internal("Error creating user", nil)
	}
	return nil
}
