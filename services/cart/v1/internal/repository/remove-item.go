package repository

import (
	"context"

	"github.com/goforj/godump"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (repo *Repository) RemoveItem(ctx context.Context, userID string, itemID int) error {
	filter := bson.M{"user_id": userID}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": itemID},
		},
	}

	collection := repo.MongoDB.Database("ecommerce").Collection("carts")

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // return value after update
		SetProjection(bson.M{"items": 1}) // field items only

	var updatedCart bson.M
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// ไม่มี cart → จบ
			return nil
		}
		return err
	}

	godump.Dump(updatedCart)

	if len(updatedCart["items"].(bson.A)) == 0 {
		if err := repo.RemoveCart(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}
