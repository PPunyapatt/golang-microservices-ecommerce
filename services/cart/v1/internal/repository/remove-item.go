package repository

import (
	"context"

	"github.com/goforj/godump"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (repo *Repository) RemoveItem(userID string, cartID, itemID int) error {
	// filter := bson.M{"user_id": userID}

	// update := bson.M{
	// 	"$pull": bson.M{
	// 		"items": bson.M{"product_id": itemID}, // ลบ item ที่ product_id ตรงกับ productID
	// 	},
	// }

	// collection := repo.MongoDB.Database("ecommerce").Collection("carts")
	// _, err := collection.UpdateOne(context.Background(), filter, update)
	// if err != nil {
	// 	return err
	// }

	// var carts bson.M
	// err = collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&carts)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		// ไม่มี cart แล้ว → return ปกติ
	// 		return nil
	// 	}
	// 	return err
	// }

	// godump.Dump(carts["items"])

	filter := bson.M{"user_id": userID}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": itemID},
		},
	}

	collection := repo.MongoDB.Database("ecommerce").Collection("carts")

	// ใช้ FindOneAndUpdate แทน UpdateOne + FindOne
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // คืนค่าหลัง update
		SetProjection(bson.M{"items": 1}) // เอาเฉพาะ field items

	// var updatedCart struct {
	// 	Items []bson.M `bson:"items"`
	// }
	var updatedCart bson.D
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// ไม่มี cart → จบ
			return nil
		}
		return err
	}

	godump.Dump(updatedCart)
	return nil
}
