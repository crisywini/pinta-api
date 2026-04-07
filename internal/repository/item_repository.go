package repository

import (
	"context"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ItemRepository struct {
	collection *mongo.Collection
}

func NewItemRepository(db *mongo.Database) *ItemRepository {
	return &ItemRepository{
		collection: db.Collection("items"),
	}
}

func (r *ItemRepository) Save(item *model.Item) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) FindAll() ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var items []model.Item

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepository) FindByID(id string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item model.Item

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) Update(id string, updated *model.Item) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"name":       updated.Name,
			"category":   updated.Category,
			"color":      updated.Color,
			"brand":      updated.Brand,
			"material":   updated.Material,
			"season":     updated.Season,
			"occasion":   updated.Occasion,
			"photo":      updated.Photo,
			"condition":  updated.Condition,
			"wear_count": updated.WearCount,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err

}
