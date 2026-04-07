package repository

import (
	"context"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OutfitRepository struct {
	collection *mongo.Collection
}

func NewOutfitRepository(db *mongo.Database) *OutfitRepository {
	return &OutfitRepository{
		collection: db.Collection("outfits"),
	}
}

func (r *OutfitRepository) Save(outfit *model.Outfit) (*model.Outfit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	outfit.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, outfit)
	if err != nil {
		return nil, err
	}

	return outfit, nil
}

func (r *OutfitRepository) FindAll() ([]model.Outfit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var outfits []model.Outfit

	if err := cursor.All(ctx, &outfits); err != nil {
		return nil, err
	}

	return outfits, nil
}

func (r *OutfitRepository) FindByID(id string) (*model.Outfit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var outfit model.Outfit
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&outfit)

	if err != nil {
		return nil, err
	}

	return &outfit, nil
}

func (r *OutfitRepository) Update(id string, updated *model.Outfit) error {
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
			"occasion":   updated.Occasion,
			"season":     updated.Season,
			"mood":       updated.Mood,
			"fragrance":  updated.Fragrance,
			"notes":      updated.Notes,
			"favorite":   updated.Favorite,
			"wear_count": updated.WearCount,
			"photo":      updated.Photo,
			"last_worn":  updated.LastWorn,
			"items":      updated.Items,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *OutfitRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
