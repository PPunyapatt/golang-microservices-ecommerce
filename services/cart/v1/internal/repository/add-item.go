package repository

import (
	"cart/v1/internal/constant"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (repo *Repository) AddItem(userID string, items []*constant.StoreItems) error {

	// items_store := []map[string]interface{}{}
	// for _, store := range items {
	// 	for _, item := range store.Items {
	// 		items_store = append(items_store, map[string]interface{}{
	// 			"store_id":     store.StoreID,
	// 			"product_id":   item.ProductID,
	// 			"prodict_name": item.ProductName,
	// 			"price":        item.Price,
	// 			"quantity":     item.Quantity,
	// 			"image_url":    item.ImageURL,
	// 		})
	// 	}
	// }

	// filter := bson.M{
	// 	"user_id": userID,
	// }

	// update := bson.M{
	// 	"$push": bson.M{
	// 		"items": bson.M{"$each": items_store},
	// 	},
	// }

	// opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	// var result bson.M
	// collection := repo.MongoDB.Database("ecommerce").Collection("carts")
	// err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	return err
	// }
	// return nil
	ctx := context.Background()
	collection := repo.MongoDB.Database("ecommerce").Collection("carts")

	items_store := []map[string]interface{}{}
	for _, store := range items {
		for _, item := range store.Items {
			itm := map[string]interface{}{
				"store_id":     store.StoreID,
				"product_id":   item.ProductID,
				"prodict_name": item.ProductName,
				"price":        item.Price,
				"quantity":     item.Quantity,
				"image_url":    item.ImageURL,
			}
			items_store = append(items_store, itm)

			filter := bson.M{
				"user_id":          userID,
				"items.product_id": item.ProductID,
			}

			updateExisting := bson.M{
				"$inc": bson.M{"items.$.quantity": item.Quantity},
				"$pull": bson.M{
					"items": bson.M{"quantity": bson.M{"$lte": 0}},
				},
			}

			res := collection.FindOneAndUpdate(ctx, filter, updateExisting)

			if res.Err() == mongo.ErrNoDocuments {
				// ถ้าไม่มีสินค้า ให้ push เข้า items หรือสร้าง cart ใหม่ถ้าไม่มี cart
				filterCart := bson.M{"user_id": userID}
				updateCart := bson.M{
					"$setOnInsert": bson.M{"user_id": userID},
					"$push":        bson.M{"items": itm},
				}
				opts := options.UpdateOne().SetUpsert(true)
				_, err := collection.UpdateOne(ctx, filterCart, updateCart, opts)
				if err != nil {
					return res.Err()
				}
			} else if res.Err() != nil {
				return res.Err()
			}
		}
	}

	return nil
}
