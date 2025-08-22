package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (repo *Repository) GetOrCreateCartByUserID(ctx context.Context, userID string, pagination *constant.Pagination) ([]*cart.CartItem, error) {
	// carts := constant.Cart{}
	// filter := bson.M{
	// 	"user_id": userID,
	// }

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$project", Value: bson.M{
			"user_id":    1,
			"create_at":  1,
			"totalItems": bson.M{"$size": "$items"},
			"items":      bson.M{"$slice": []interface{}{"$items", pagination.Offset, pagination.Limit}},
		}}},
	}

	collection := repo.MongoDB.Database("ecommerce").Collection("carts")

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var carts []constant.Cart
	if err := cursor.All(ctx, &carts); err != nil {
		return nil, err
	}

	log.Println("Carts: ", carts)

	pipeline_count := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$project", Value: bson.M{
			"totalItems": bson.M{"$size": "$items"},
		}}},
	}

	items := []*cart.CartItem{}
	if len(carts) > 0 {
		for _, item := range carts[0].Items {
			items = append(items, &cart.CartItem{
				ProductId:   int32(item.ProductID),
				ProductName: item.ProductName,
				Quantity:    int32(item.Quantity),
				Price:       item.Price,
				ImageUrl:    item.ImageURL,
			})
		}

		log.Println("items: ", items)
		// ---------------- Find total count items ----------------
		cursor, err = collection.Aggregate(ctx, pipeline_count)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}

		pagination.TotalCount = int32(len(results))
	} else {
		pagination.TotalCount = int32(0)
	}

	// err = collection.FindOne(ctx, filter).Decode(&cart)
	// if err != nil {
	// 	return nil, err
	// }

	// if err != mongo.ErrNoDocuments {
	// 	return nil, err
	// }

	return items, nil
}
