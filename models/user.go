package models

import (
	"context"
	"fmt"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	ForgetToken  string             `json:"forget_token,omitempty" bson:"forget_token,omitempty"`
	Verified     bool               `json:"verified,omitempty" bson:"verified,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func CheckEmailTaken(email string) bool {
	db := mongodb.GetMongoDBConnection()
	match := db.Collection("Users").FindOne(context.TODO(), bson.M{
		"email": email,
	})
	return match.Err() == nil
}

func (n *User) GetAllUsers(c context.Context) ([]User, error) {
	db := mongodb.GetMongoDBConnection()
	cursor, err := db.Collection("Users").Find(context.TODO(), bson.M{})
	if err != nil {
		return []User{}, err
	}

	result := []User{}

	for cursor.Next(context.TODO()) {
		elem := User{}
		cursor.Decode(&elem)
		result = append(result, elem)
	}

	return result, nil

}

func (n *User) Register() error {
	if CheckEmailTaken(n.Email) {
		return fmt.Errorf("register.invalid.email")
	}

	db := mongodb.GetMongoDBConnection()

	n.ID = primitive.NewObjectID()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	_, err := db.Collection("Users").InsertOne(context.TODO(), n)
	if err != nil {
		return err
	}

	return nil
}

func (n *User) Login(email string, password string) (User, error) {
	db := mongodb.GetMongoDBConnection()
	match := db.Collection("Users").FindOne(context.TODO(), bson.M{
		"email":    email,
		"password": password,
	})

	if match.Err() != nil {
		return User{}, fmt.Errorf("register.invalid.email")
	}

	var user User
	match.Decode(&user)

	return user, nil
}
