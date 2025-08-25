package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (repo *Repository) GetOrCreateCartByUserID(ctx context.Context, userID string, pagination *constant.Pagination) ([]*cart.StoreItems, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$unwind", Value: "$items"}},
		{{Key: "$skip", Value: pagination.Offset}}, // pagination offset
		{{Key: "$limit", Value: pagination.Limit}}, // pagination limit
		{{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"user_id":  "$user_id",
					"store_id": "$items.store_id",
				},
				"items": bson.M{"$push": "$items"},
			},
		}},
		{{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"user_id": "$_id.user_id",
				},
				"store_items": bson.M{
					"$push": bson.M{
						"store_id": "$_id.store_id",
						"items":    "$items",
					},
				},
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"_id":         0,
				"user_id":     "$_id.user_id",
				"store_items": 1,
			},
		}},
	}

	collection := repo.MongoDB.Database("ecommerce").Collection("carts")

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var carts []*constant.Cart

	if err := cursor.All(ctx, &carts); err != nil {
		return nil, err
	}

	// godump.Dump(carts)

	// ------------------- Count -------------------
	countPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$unwind", Value: "$items"}},
		{{Key: "$count", Value: "total_count"}},
	}

	cursorCount, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursorCount.All(ctx, &results); err != nil {
		return nil, err
	}
	// godump.Dump(results)

	storeItems := []*cart.StoreItems{}
	if carts != nil {
		for _, store := range carts[0].StoreItems {
			items := []*cart.CartItem{}
			for _, item := range store.Items {
				items = append(items, &cart.CartItem{
					ProductId:   int32(item.ProductID),
					ProductName: item.ProductName,
					Price:       item.Price,
					Quantity:    int32(item.Quantity),
					ImageUrl:    item.ImageURL,
				})
			}
			storeItems = append(storeItems, &cart.StoreItems{
				StoreID: int32(store.StoreID),
				Items:   items,
			})
		}

		pagination.TotalCount = int32(results[0]["total_count"].(int32))
	} else {
		pagination.TotalCount = 0
	}

	return storeItems, nil
}
