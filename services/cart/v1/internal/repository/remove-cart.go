package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (repo *Repository) RemoveCart(ctx context.Context, userID string) error {
	collection := repo.MongoDB.Database("ecommerce").Collection("carts")
	filter := bson.M{"user_id": userID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
