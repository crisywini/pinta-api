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
