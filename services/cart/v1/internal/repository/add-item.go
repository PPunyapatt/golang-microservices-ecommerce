package repository

import (
	"cart/v1/internal/constant"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (repo *Repository) AddItem(userID string, items []*constant.Item) error {

	// cart, err := repo.GetOrCreateCartByUserID(userID)
	// if err != nil {
	// 	return err
	// }

	// for _, item := range items {
	// 	item.CartID = cart.CartID
	// }
	// // godump.Dump(items)
	// if err := repo.GormDB.Omit("updated_at").Create(&items); err != nil {
	// 	return err.Error
	// }

	filter := bson.M{
		"user_id": userID,
	}

	update := bson.M{
		"$push": bson.M{
			"items": bson.M{"$each": items},
		},
		"$setOnInsert": bson.M{
			"userId":    "USER123",
			"createdAt": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result bson.M
	collection := repo.MongoDB.Database("ecommerce").Collection("carts")
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}
