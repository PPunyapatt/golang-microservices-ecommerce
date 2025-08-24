package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (repo *Repository) RemoveItem(userID string, cartID, itemID int) error {
	filter := bson.M{"user_id": userID}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": itemID}, // ลบ item ที่ product_id ตรงกับ productID
		},
	}

	collection := repo.MongoDB.Database("ecommerce").Collection("carts")
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
