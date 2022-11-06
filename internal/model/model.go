/*
Package model containing structs and functions
for database transaction
*/
package model

import (
	"context"
	"fmt"

	"github.com/reyhanfikridz/ecom-order-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Order contain order detail
type Order struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" form:"_id,omitempty"`
	OrderNumber         string             `bson:"order_number" json:"order_number" form:"order_number"`
	Status              string             `bson:"status" json:"status" form:"status"`
	Qty                 int                `bson:"qty" json:"qty" form:"qty"`
	TotalPrice          float64            `bson:"total_price" json:"total_price" form:"total_price"`
	BuyerID             int                `bson:"buyer_id" json:"buyer_id" form:"buyer_id"`
	BuyerFullName       string             `bson:"buyer_full_name" json:"buyer_full_name" form:"buyer_full_name"`
	BuyerAddress        string             `bson:"buyer_address" json:"buyer_address" form:"buyer_address"`
	ProductID           int                `bson:"product_id" json:"product_id" form:"product_id"`
	ProductSKU          string             `bson:"product_sku" json:"product_sku" form:"product_sku"`
	ProductName         string             `bson:"product_name" json:"product_name" form:"product_name"`
	ProductPrice        float64            `bson:"product_price" json:"product_price" form:"product_price"`
	ProductWeight       float32            `bson:"product_weight" json:"product_weight" form:"product_weight"`
	ProductDescription  string             `bson:"product_description" json:"product_description" form:"product_description"`
	ProductStock        int                `bson:"product_stock" json:"product_stock" form:"product_stock"`
	ProductUserID       int                `bson:"product_user_id" json:"product_user_id" form:"product_user_id"`
	ProductUserFullName string             `bson:"product_user_full_name" json:"product_user_full_name" form:"product_user_full_name"`
	ProductImagesPath   []string           `bson:"product_images_path" json:"product_images_path" form:"product_images_path"`
}

// InsertOrder insert order document to orders collection
func InsertOrder(ctx context.Context, oc *mongo.Collection, o Order) (Order, error) {
	// get random order number until new one found
	// by checking it in collection
	var orderNumber string
	for {
		orderNumber = utils.GetRandomOrderNumber()

		var tmp bson.M
		err := oc.FindOne(ctx, bson.D{
			primitive.E{Key: "order_number", Value: orderNumber},
		}).Decode(&tmp)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
		}
	}

	// insert order to database
	o.OrderNumber = orderNumber
	result, err := oc.InsertOne(ctx, o)
	if err != nil {
		return o, err
	}

	// get order ID
	o.ID = result.InsertedID.(primitive.ObjectID)

	return o, nil
}

// GetOrder get order document by some key from orders collection
func GetOrder(ctx context.Context, oc *mongo.Collection,
	filter bson.M) (Order, error) {
	o := Order{}

	err := oc.FindOne(ctx, filter).Decode(&o)
	if err != nil {
		return o, err
	}

	return o, nil
}

// UpdateOrder update order document by some key in orders collection
func UpdateOrder(ctx context.Context, oc *mongo.Collection,
	filter bson.M, oUpdate Order) error {
	// set fields that need to be updated
	value := bson.M{}
	if oUpdate.Status != "" {
		value["status"] = oUpdate.Status
	}
	if oUpdate.Qty != 0 {
		value["qty"] = oUpdate.Qty
	}
	if oUpdate.TotalPrice != 0 {
		value["total_price"] = oUpdate.TotalPrice
	}
	if oUpdate.BuyerID != 0 {
		value["buyer_id"] = oUpdate.BuyerID
	}
	if oUpdate.BuyerFullName != "" {
		value["buyer_full_name"] = oUpdate.BuyerFullName
	}
	if oUpdate.BuyerAddress != "" {
		value["buyer_address"] = oUpdate.BuyerAddress
	}
	if oUpdate.ProductID != 0 {
		value["product_id"] = oUpdate.ProductID
	}
	if oUpdate.ProductSKU != "" {
		value["product_sku"] = oUpdate.ProductSKU
	}
	if oUpdate.ProductName != "" {
		value["product_name"] = oUpdate.ProductName
	}
	if oUpdate.ProductPrice != 0 {
		value["product_price"] = oUpdate.ProductPrice
	}
	if oUpdate.ProductWeight != 0 {
		value["product_weight"] = oUpdate.ProductWeight
	}
	if oUpdate.ProductDescription != "" {
		value["product_description"] = oUpdate.ProductDescription
	}
	if oUpdate.ProductStock != 0 {
		value["product_stock"] = oUpdate.ProductStock
	}
	if oUpdate.ProductUserID != 0 {
		value["product_user_id"] = oUpdate.ProductUserID
	}
	if oUpdate.ProductUserFullName != "" {
		value["product_user_full_name"] = oUpdate.ProductUserFullName
	}
	if oUpdate.ProductImagesPath != nil {
		value["product_images_path"] = oUpdate.ProductImagesPath
	}

	fields := bson.M{"$set": value}

	// update order
	result, err := oc.UpdateOne(ctx, filter, fields)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no data updated")
	}

	return nil
}

// DeleteOrder delete order document by some key in orders collection
func DeleteOrder(ctx context.Context, oc *mongo.Collection,
	filter bson.M) error {
	// delete order
	_, err := oc.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetOrders get order documents by some key in orders collection
func GetOrders(ctx context.Context, oc *mongo.Collection,
	filter bson.M) ([]Order, error) {
	orders := []Order{}

	// get orders cursor
	cur, err := oc.Find(ctx, filter)
	if err != nil {
		return orders, err
	}
	defer cur.Close(ctx)

	// loop orders in cursor and put it into slice of order
	for cur.Next(ctx) {
		var o Order
		err = cur.Decode(&o)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
