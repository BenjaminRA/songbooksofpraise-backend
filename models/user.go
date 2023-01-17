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
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Admin     bool               `json:"admin" bson:"admin"`
	Editor    bool               `json:"editor" bson:"editor"`
	Moderator bool               `json:"moderator" bson:"moderator"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func CheckEmailTaken(email string) bool {
	db := mongodb.GetMongoDBConnection()
	match := db.Collection("Users").FindOne(context.TODO(), bson.M{
		"email": email,
	})
	return match.Err() == nil
}

func CheckEmailTakenWithId(email string, user_id primitive.ObjectID) bool {
	db := mongodb.GetMongoDBConnection()
	match := db.Collection("Users").FindOne(context.TODO(), bson.M{
		"_id":   bson.M{"$ne": user_id},
		"email": email,
	})
	return match.Err() == nil
}

func (n *User) GetUserById(user_id string) (User, error) {
	db := mongodb.GetMongoDBConnection()
	objectID, _ := primitive.ObjectIDFromHex(user_id)

	match := db.Collection("Users").FindOne(context.TODO(), bson.M{
		"_id": objectID,
	})

	if match.Err() != nil {
		return User{}, match.Err()
	}

	var user User
	match.Decode(&user)

	return user, nil
}

func (n *User) GetAllUsers() ([]User, error) {
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
	n.Admin = false
	n.Editor = true
	n.Moderator = false
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

func (n *User) UpdateUser() error {
	if CheckEmailTakenWithId(n.Email, n.ID) {
		return fmt.Errorf("register.invalid.email")
	}

	db := mongodb.GetMongoDBConnection()
	_, err := db.Collection("Users").UpdateOne(context.TODO(), bson.M{
		"email": n.Email,
	}, bson.M{
		"$set": bson.M{
			"first_name": n.FirstName,
			"last_name":  n.LastName,
			"email":      n.Email,
			"admin":      n.Admin,
			"moderator":  n.Moderator,
			"editor":     n.Editor,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *User) DeleteUser() error {
	db := mongodb.GetMongoDBConnection()
	_, err := db.Collection("Users").DeleteOne(context.TODO(), bson.M{
		"email": n.Email,
	})
	if err != nil {
		return err
	}

	return nil
}
