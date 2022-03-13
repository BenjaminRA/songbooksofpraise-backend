package models

import (
	"context"
	"fmt"
	"sync"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID         primitive.ObjectID   `json:"_id" bson:"_id"`
	All        bool                 `json:"all" bson:"all"`
	Category   string               `json:"category" bson:"category"`
	SongbookID primitive.ObjectID   `json:"songbook_id" bson:"songbook_id"`
	ParentID   primitive.ObjectID   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	ChildrenID []primitive.ObjectID `json:"children_id,omitempty" bson:"children_id,omitempty"`
	Children   []Category           `json:"children,omitempty" bson:"children,omitempty"`
	Songs      []Song               `json:"songs,omitempty" bson:"songs,omitempty"`
}

func (n *Category) GetAllCategories() []Category {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Categories").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"parent_id": primitive.Null{}}},
	})
	if err != nil {
		panic(err)
	}

	result := []Category{}

	var wg sync.WaitGroup
	for cursor.Next(context.TODO()) {
		elem := Category{}
		cursor.Decode(&elem)
		fmt.Println(elem.Category, elem.All)
		if elem.ID.Hex() != "000000000000000000000000" {
			wg.Add(1)
			elem.Children = elem.GetChildren()
			wg.Done()
		}

		result = append(result, elem)
	}
	wg.Wait()

	return result
}

func (n *Category) GetCategoryById(id string) Category {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	cursor, err := db.Collection("Categories").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"_id": object_id}},
		{"$lookup": bson.M{
			"from":         "Songs",
			"localField":   "_id",
			"foreignField": "categories_id",
			"pipeline": []bson.M{
				{
					"$sort": bson.M{
						"number": 1,
					},
				},
			},
			"as": "songs",
		}},
	})
	if err != nil {
		panic(err)
	}

	result := []Category{}

	for cursor.Next(context.TODO()) {
		elem := Category{}
		cursor.Decode(&elem)

		if elem.ID.Hex() != "000000000000000000000000" {
			elem.Children = elem.GetChildren()
		}

		result = append(result, elem)
	}

	return result[0]
}

func (n *Category) GetChildren() []Category {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Categories").Find(context.TODO(), bson.M{
		"parent_id": n.ID,
	})
	if err != nil {
		panic(err)
	}

	result := []Category{}

	wg := sync.WaitGroup{}

	for cursor.Next(context.TODO()) {
		wg.Add(1)

		elem := Category{}
		cursor.Decode(&elem)

		if elem.ID.Hex() != "000000000000000000000000" {
			elem.Children = elem.GetChildren()
		}

		result = append(result, elem)

		wg.Done()
	}

	wg.Wait()

	return result
}

func (n *Category) CreateCategory() (Category, error) {
	db := mongodb.GetMongoDBConnection()
	n.ID = primitive.NewObjectID()

	_, err := db.Collection("Categories").InsertOne(context.TODO(), n)
	if err != nil {
		return Category{}, err
	}

	return new(Category).GetCategoryById(n.ID.Hex()), nil
}

func (n *Category) UpdateCategory() error {
	db := mongodb.GetMongoDBConnection()

	update := bson.M{
		"$set": bson.M{
			"category": n.Category,
		},
	}
	if n.ParentID.Hex() == "000000000000000000000000" {
		update["$unset"] = bson.M{
			"parent_id": "",
		}
	} else {
		update["$set"].(bson.M)["parent_id"] = n.ParentID
	}

	_, err := db.Collection("Categories").UpdateOne(context.TODO(), bson.M{
		"_id": n.ID,
	}, update)
	if err != nil {
		return (err)
	}

	return nil
}

func (n *Category) DeleteCategory() error {
	db := mongodb.GetMongoDBConnection()

	n.Children = n.GetChildren()

	if len(n.Children) > 0 {
		n.deleteAllChildren()
	}

	_, err := db.Collection("Categories").DeleteOne(context.TODO(), bson.M{
		"_id": n.ID,
	})

	if err != nil {
		return (err)
	}

	// Deleting category from songs
	_, err = db.Collection("Songs").UpdateMany(context.TODO(), bson.M{
		"categories_id": n.ID,
	}, bson.M{
		"$pull": bson.M{
			"categories_id": n.ID,
		},
	})

	if err != nil {
		return (err)
	}

	return nil
}

func (n *Category) deleteAllChildren() {
	for _, child := range n.Children {
		child.Children = child.GetChildren()
		if len(child.Children) > 0 {
			child.deleteAllChildren()
		}

		child.DeleteCategory()
	}
}

func AllToFirst(t *[]Category) {
	todo_idx := -1

	for i, value := range *t {
		if value.All {
			todo_idx = i
			break
		}
	}

	if todo_idx != -1 {
		for i := todo_idx; i > 0; i-- {
			(*t)[i], (*t)[i-1] = (*t)[i-1], (*t)[i]
		}
	}
}
