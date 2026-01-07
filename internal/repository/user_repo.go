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

func (ur *UserRepository) UpdatePasswordToken(ctx context.Context, user *models.User) error {
	filter := bson.M{"_id": user.ID, "email": user.Email}
	updates := bson.M{
		"$set": bson.M{
			"password_reset_token":         user.Password_reset_token,
			"password_reset_token_expires": user.Password_reset_token_expires,
		},
	}

	_, err := ur.Collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) UpdatePassword(ctx context.Context, token string, req *models.UpdatePasswordRequest) error {

	filter := bson.M{"password_reset_token": token, "password_reset_token_expires": bson.M{"$gt": time.Now()}}
	updates := bson.M{
		"$set": bson.M{"password": req.Password, "updated_at": time.Now()},
		"$unset": bson.M{
			"password_reset_token":         "",
			"password_reset_token_expires": "",
		},
	}

	result, err := ur.Collection.UpdateOne(ctx, filter, updates)
	if err != nil || result.MatchedCount == 0 {
		fmt.Println("error:", err)
		fmt.Println("result:", result)
		return utils.BadRequest("Invalid or expired token", nil)
	}
	return nil

}
