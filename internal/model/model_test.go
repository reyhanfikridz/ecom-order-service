/*
Package model containing structs and functions
for database transaction
*/
package model

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/reyhanfikridz/ecom-order-service/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestMain do some test before and after all testing in the package
func TestMain(m *testing.M) {
	ctx := context.Background()

	// init all config before can be used
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("There's an error when initialize config => %s", err)
	}

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		log.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// remove all data in all collection before run all test
	for _, collection := range collections {
		_, err = collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections before run all test => %s", err)
		}
	}

	// run all test
	m.Run()

	// remove all data in all collection before run all test
	for _, collection := range collections {
		_, err = collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// TestCreateOrder test CreateOrder
func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		t.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// create order struct
	o := Order{
		Status:              "in-cart",
		Qty:                 2,
		TotalPrice:          2000001,
		BuyerID:             1,
		BuyerFullName:       "George Marcus",
		BuyerAddress:        "Buyer Street",
		ProductID:           1,
		ProductSKU:          "testsku",
		ProductName:         "product name",
		ProductPrice:        1000000.50,
		ProductWeight:       1.5,
		ProductDescription:  "product description",
		ProductStock:        100,
		ProductUserID:       10,
		ProductUserFullName: "Reyhan",
		ProductImagesPath:   []string{"product 1.1.jpg", "product 1.2.jpg"},
	}

	// test insert order success
	o, err = InsertOrder(ctx, collections["orders"], o)
	if err != nil {
		t.Fatalf("Expected insert success, but got error => %s", err)
	}
	if o.ID == primitive.NilObjectID {
		t.Errorf("Expected ID not nil, but got nil")
	}

	// remove all data order after test
	_, err = collections["orders"].DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("There's an error when truncating "+
			"orders collection after test => %s", err)
	}
}

// TestGetOrder test GetOrder
//
// Required for the test: CreateOrder
func TestGetOrder(t *testing.T) {
	ctx := context.Background()

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		t.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// create order struct
	o := Order{
		Status:              "in-cart",
		Qty:                 2,
		TotalPrice:          2000001,
		BuyerID:             1,
		BuyerFullName:       "George Marcus",
		BuyerAddress:        "Buyer Street",
		ProductID:           1,
		ProductSKU:          "testsku",
		ProductName:         "product name",
		ProductPrice:        1000000.50,
		ProductWeight:       1.5,
		ProductDescription:  "product description",
		ProductStock:        100,
		ProductUserID:       10,
		ProductUserFullName: "Reyhan",
		ProductImagesPath:   []string{"product 1.1.jpg", "product 1.2.jpg"},
	}

	// insert order
	o, err = InsertOrder(ctx, collections["orders"], o)
	if err != nil {
		t.Fatalf("There's an error when inserting order data => %s", err)
	}

	// test get order by order number and check the result
	filter := bson.M{"order_number": o.OrderNumber}
	result, err := GetOrder(ctx, collections["orders"], filter)
	if err != nil {
		t.Errorf("Expected get order by order number success, "+
			"but got error => %s", err)
	}
	if o.ID != result.ID {
		t.Errorf("Expected ID %s, but got %s", o.ID, result.ID)
	}
	if o.OrderNumber != result.OrderNumber {
		t.Errorf("Expected OrderNumber %s, but OrderNumber %s",
			o.OrderNumber, result.OrderNumber)
	}
	if o.Status != result.Status {
		t.Errorf("Expected Status %s, but Status %s",
			o.Status, result.Status)
	}
	if o.Qty != result.Qty {
		t.Errorf("Expected Qty %d, but Qty %d",
			o.Qty, result.Qty)
	}
	if o.TotalPrice != result.TotalPrice {
		t.Errorf("Expected TotalPrice %f, but TotalPrice %f",
			o.TotalPrice, result.TotalPrice)
	}
	if o.BuyerID != result.BuyerID {
		t.Errorf("Expected BuyerID %d, but BuyerID %d",
			o.BuyerID, result.BuyerID)
	}
	if o.BuyerFullName != result.BuyerFullName {
		t.Errorf("Expected BuyerFullName %s, but BuyerFullName %s",
			o.BuyerFullName, result.BuyerFullName)
	}
	if o.BuyerAddress != result.BuyerAddress {
		t.Errorf("Expected BuyerAddress %s, but BuyerAddress %s",
			o.BuyerAddress, result.BuyerAddress)
	}
	if o.ProductID != result.ProductID {
		t.Errorf("Expected ProductID %d, but ProductID %d",
			o.ProductID, result.ProductID)
	}
	if o.ProductSKU != result.ProductSKU {
		t.Errorf("Expected ProductSKU %s, but ProductSKU %s",
			o.ProductSKU, result.ProductSKU)
	}
	if o.ProductName != result.ProductName {
		t.Errorf("Expected ProductName %s, but ProductName %s",
			o.ProductName, result.ProductName)
	}
	if o.ProductPrice != result.ProductPrice {
		t.Errorf("Expected ProductPrice %f, but ProductPrice %f",
			o.ProductPrice, result.ProductPrice)
	}
	if o.ProductWeight != result.ProductWeight {
		t.Errorf("Expected ProductWeight %f, but ProductWeight %f",
			o.ProductWeight, result.ProductWeight)
	}
	if o.ProductDescription != result.ProductDescription {
		t.Errorf("Expected ProductDescription %s, but ProductDescription %s",
			o.ProductDescription, result.ProductDescription)
	}
	if o.ProductStock != result.ProductStock {
		t.Errorf("Expected ProductStock %d, but ProductStock %d",
			o.ProductStock, result.ProductStock)
	}
	if o.ProductUserID != result.ProductUserID {
		t.Errorf("Expected ProductUserID %d, but ProductUserID %d",
			o.ProductUserID, result.ProductUserID)
	}
	if o.ProductUserFullName != result.ProductUserFullName {
		t.Errorf("Expected ProductUserFullName %s, but ProductUserFullName %s",
			o.ProductUserFullName, result.ProductUserFullName)
	}
	if len(o.ProductImagesPath) != len(result.ProductImagesPath) {
		t.Errorf("Expected total ProductImagesPath %d, but total ProductImagesPath %d",
			len(o.ProductImagesPath), len(result.ProductImagesPath))
	} else {
		for _, expectedProductImagePath := range o.ProductImagesPath {
			imageExist := false
			for _, resultProductImagePath := range result.ProductImagesPath {
				if expectedProductImagePath == resultProductImagePath {
					imageExist = true
				}
			}
			if !imageExist {
				t.Errorf("Expected ProductImagesPath %s not found",
					expectedProductImagePath)
			}
		}
	}

	// remove all data order after test
	_, err = collections["orders"].DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("There's an error when truncating "+
			"orders collection after test => %s", err)
	}
}

// TestUpdateOrder test UpdateOrder
//
// Required for the test:
//
// - CreateOrder
//
// - GetOrder
func TestUpdateOrder(t *testing.T) {
	ctx := context.Background()

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		t.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// create order struct
	o := Order{
		Status:              "in-cart",
		Qty:                 2,
		TotalPrice:          2000001,
		BuyerID:             1,
		BuyerFullName:       "George Marcus",
		BuyerAddress:        "Buyer Street",
		ProductID:           1,
		ProductSKU:          "testsku",
		ProductName:         "product name",
		ProductPrice:        1000000.50,
		ProductWeight:       1.5,
		ProductDescription:  "product description",
		ProductStock:        100,
		ProductUserID:       10,
		ProductUserFullName: "Reyhan",
		ProductImagesPath:   []string{"product 1.1.jpg", "product 1.2.jpg"},
	}

	// insert order
	o, err = InsertOrder(ctx, collections["orders"], o)
	if err != nil {
		t.Fatalf("There's an error when inserting order data => %s", err)
	}

	// test update order and check the result
	oUpdate := Order{
		OrderNumber:         o.OrderNumber,
		Status:              "processed",
		Qty:                 3,
		TotalPrice:          6000001.50,
		BuyerID:             2,
		BuyerFullName:       "Sean Marco",
		BuyerAddress:        "Buyer Street 2",
		ProductID:           2,
		ProductSKU:          "testsku2",
		ProductName:         "product name 2",
		ProductPrice:        2000000.50,
		ProductWeight:       2.5,
		ProductDescription:  "product description 2",
		ProductStock:        200,
		ProductUserID:       20,
		ProductUserFullName: "Fikri",
		ProductImagesPath:   []string{"product 2.1.jpg", "product 2.2.jpg", "product 2.3.jpg"},
	}

	filter := bson.M{"order_number": oUpdate.OrderNumber}
	err = UpdateOrder(ctx, collections["orders"], filter, oUpdate)
	if err != nil {
		t.Errorf("Expected update order by order number success, "+
			"but got error => %s", err)
	}

	filter = bson.M{"order_number": o.OrderNumber}
	result, err := GetOrder(ctx, collections["orders"], filter)
	if err != nil {
		t.Errorf("There's an error when get order by order number => %s",
			err)
	}
	if oUpdate.OrderNumber != result.OrderNumber {
		t.Errorf("Expected OrderNumber %s, but OrderNumber %s",
			oUpdate.OrderNumber, result.OrderNumber)
	}
	if oUpdate.Status != result.Status {
		t.Errorf("Expected Status %s, but Status %s",
			oUpdate.Status, result.Status)
	}
	if oUpdate.Qty != result.Qty {
		t.Errorf("Expected Qty %d, but Qty %d",
			oUpdate.Qty, result.Qty)
	}
	if oUpdate.TotalPrice != result.TotalPrice {
		t.Errorf("Expected TotalPrice %f, but TotalPrice %f",
			oUpdate.TotalPrice, result.TotalPrice)
	}
	if oUpdate.BuyerID != result.BuyerID {
		t.Errorf("Expected BuyerID %d, but BuyerID %d",
			oUpdate.BuyerID, result.BuyerID)
	}
	if oUpdate.BuyerFullName != result.BuyerFullName {
		t.Errorf("Expected BuyerFullName %s, but BuyerFullName %s",
			oUpdate.BuyerFullName, result.BuyerFullName)
	}
	if oUpdate.BuyerAddress != result.BuyerAddress {
		t.Errorf("Expected BuyerAddress %s, but BuyerAddress %s",
			oUpdate.BuyerAddress, result.BuyerAddress)
	}
	if oUpdate.ProductID != result.ProductID {
		t.Errorf("Expected ProductID %d, but ProductID %d",
			oUpdate.ProductID, result.ProductID)
	}
	if oUpdate.ProductSKU != result.ProductSKU {
		t.Errorf("Expected ProductSKU %s, but ProductSKU %s",
			oUpdate.ProductSKU, result.ProductSKU)
	}
	if oUpdate.ProductName != result.ProductName {
		t.Errorf("Expected ProductName %s, but ProductName %s",
			oUpdate.ProductName, result.ProductName)
	}
	if oUpdate.ProductPrice != result.ProductPrice {
		t.Errorf("Expected ProductPrice %f, but ProductPrice %f",
			oUpdate.ProductPrice, result.ProductPrice)
	}
	if oUpdate.ProductWeight != result.ProductWeight {
		t.Errorf("Expected ProductWeight %f, but ProductWeight %f",
			oUpdate.ProductWeight, result.ProductWeight)
	}
	if oUpdate.ProductDescription != result.ProductDescription {
		t.Errorf("Expected ProductDescription %s, but ProductDescription %s",
			oUpdate.ProductDescription, result.ProductDescription)
	}
	if oUpdate.ProductStock != result.ProductStock {
		t.Errorf("Expected ProductStock %d, but ProductStock %d",
			oUpdate.ProductStock, result.ProductStock)
	}
	if oUpdate.ProductUserID != result.ProductUserID {
		t.Errorf("Expected ProductUserID %d, but ProductUserID %d",
			oUpdate.ProductUserID, result.ProductUserID)
	}
	if oUpdate.ProductUserFullName != result.ProductUserFullName {
		t.Errorf("Expected ProductUserFullName %s, but ProductUserFullName %s",
			oUpdate.ProductUserFullName, result.ProductUserFullName)
	}
	if len(oUpdate.ProductImagesPath) != len(result.ProductImagesPath) {
		t.Errorf("Expected total ProductImagesPath %d, but total ProductImagesPath %d",
			len(oUpdate.ProductImagesPath), len(result.ProductImagesPath))
	} else {
		for _, expectedProductImagePath := range oUpdate.ProductImagesPath {
			imageExist := false
			for _, resultProductImagePath := range result.ProductImagesPath {
				if expectedProductImagePath == resultProductImagePath {
					imageExist = true
				}
			}
			if !imageExist {
				t.Errorf("Expected ProductImagesPath %s not found",
					expectedProductImagePath)
			}
		}
	}

	// test update order no data updated
	oUpdate = Order{
		OrderNumber: "this order number not exist in database",
	}

	filter = bson.M{"order_number": oUpdate.OrderNumber}
	err = UpdateOrder(ctx, collections["orders"], filter, oUpdate)
	if err == nil {
		t.Errorf("Expected error no data updated, but got no error")
	} else {
		if err.Error() != "no data updated" {
			t.Errorf("Expected error no data updated, but error %s", err)
		}
	}

	// remove all data order after test
	_, err = collections["orders"].DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("There's an error when truncating "+
			"orders collection after test => %s", err)
	}
}

// TestDeleteOrder test DeleteOrder
//
// Required for the test:
//
// - CreateOrder
//
// - GetOrder
func TestDeleteOrder(t *testing.T) {
	ctx := context.Background()

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		t.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// create order struct
	o := Order{
		Status:              "in-cart",
		Qty:                 2,
		TotalPrice:          2000001,
		BuyerID:             1,
		BuyerFullName:       "George Marcus",
		BuyerAddress:        "Buyer Street",
		ProductID:           1,
		ProductSKU:          "testsku",
		ProductName:         "product name",
		ProductPrice:        1000000.50,
		ProductWeight:       1.5,
		ProductDescription:  "product description",
		ProductStock:        100,
		ProductUserID:       10,
		ProductUserFullName: "Reyhan",
		ProductImagesPath:   []string{"product 1.1.jpg", "product 1.2.jpg"},
	}

	// insert order
	o, err = InsertOrder(ctx, collections["orders"], o)
	if err != nil {
		t.Fatalf("There's an error when inserting order data => %s", err)
	}

	// test delete order by order number and check the result
	filter := bson.M{"order_number": o.OrderNumber}

	err = DeleteOrder(ctx, collections["orders"], filter)
	if err != nil {
		t.Errorf("Expected delete order by order number success, "+
			"but got error => %s", err)
	}

	_, err = GetOrder(ctx, collections["orders"], filter)
	if err == nil {
		t.Errorf("Expected error errNoDocuments, but got no error")
	} else {
		if err != mongo.ErrNoDocuments {
			t.Errorf("Expected error errNoDocuments, but got error %s", err)
		}
	}

	// remove all data order after test
	_, err = collections["orders"].DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("There's an error when truncating "+
			"orders collection after test => %s", err)
	}
}

// TestGetOrders test GetOrders
//
// Required for the test: CreateOrder
func TestGetOrders(t *testing.T) {
	ctx := context.Background()

	// get map of collection
	collections, err := getTestingCollections(ctx)
	if err != nil {
		t.Fatalf("There's an error when getting "+
			"mongodb collections => %s", err)
	}

	// create orders struct
	orders := []Order{
		{
			Status:              "in-cart",
			Qty:                 2,
			TotalPrice:          2000001,
			BuyerID:             1,
			BuyerFullName:       "George Marcus",
			BuyerAddress:        "Buyer Street",
			ProductID:           1,
			ProductSKU:          "testsku",
			ProductName:         "product name",
			ProductPrice:        1000000.50,
			ProductWeight:       1.5,
			ProductDescription:  "product description",
			ProductStock:        100,
			ProductUserID:       10,
			ProductUserFullName: "Reyhan",
			ProductImagesPath:   []string{"product 1.1.jpg", "product 1.2.jpg"},
		},
		{
			Status:              "done",
			Qty:                 2,
			TotalPrice:          4000001,
			BuyerID:             1,
			BuyerFullName:       "George Marcus",
			BuyerAddress:        "Buyer Street",
			ProductID:           2,
			ProductSKU:          "testsku2",
			ProductName:         "product name 2",
			ProductPrice:        2000000.50,
			ProductWeight:       2.5,
			ProductDescription:  "product description 2",
			ProductStock:        200,
			ProductUserID:       20,
			ProductUserFullName: "Fikri",
			ProductImagesPath:   []string{"product 2.1.jpg", "product 2.2.jpg", "product 2.3.jpg"},
		},
		{
			Status:              "in-cart",
			Qty:                 2,
			TotalPrice:          6000001,
			BuyerID:             2,
			BuyerFullName:       "Linda",
			BuyerAddress:        "Buyer Street 2",
			ProductID:           3,
			ProductSKU:          "testsku 3",
			ProductName:         "product name 3",
			ProductPrice:        3000000.50,
			ProductWeight:       3.5,
			ProductDescription:  "product description 3",
			ProductStock:        300,
			ProductUserID:       30,
			ProductUserFullName: "Dzikriansyah",
			ProductImagesPath:   []string{},
		},
	}

	// insert orders
	for i, o := range orders {
		orders[i], err = InsertOrder(ctx, collections["orders"], o)
		if err != nil {
			t.Fatalf("There's an error when inserting order data => %s", err)
		}
	}

	// create testing table
	testTable := []struct {
		TestName        string
		Filter          bson.M
		ExpectedResults []Order
	}{
		{
			TestName:        "Test Get All Order",
			Filter:          bson.M{},
			ExpectedResults: orders,
		},
		{
			TestName:        "Test Get Order By Buyer ID <1>",
			Filter:          bson.M{"buyer_id": 1},
			ExpectedResults: []Order{orders[0], orders[1]},
		},
		{
			TestName:        "Test Get Order By Buyer ID <1> and Status <in-cart>",
			Filter:          bson.M{"buyer_id": 1, "status": "in-cart"},
			ExpectedResults: []Order{orders[0]},
		},
	}

	// do test in test table
	for _, test := range testTable {
		// test get orders and check the result
		results, err := GetOrders(ctx, collections["orders"], test.Filter)
		if err != nil {
			t.Errorf("[%s] Expected get orders success, "+
				"but got error => %s", test.TestName, err)
		}

		if len(test.ExpectedResults) != len(results) {
			t.Errorf("[%s] Expected total order %d, but got %d", test.TestName,
				len(test.ExpectedResults), len(results))
		}

		for _, expectedResult := range test.ExpectedResults {
			resultExist := false
			for _, result := range results {
				if expectedResult.ID == result.ID &&
					expectedResult.OrderNumber == result.OrderNumber &&
					expectedResult.Status == result.Status &&
					expectedResult.Qty == result.Qty &&
					expectedResult.TotalPrice == result.TotalPrice &&
					expectedResult.BuyerID == result.BuyerID &&
					expectedResult.BuyerFullName == result.BuyerFullName &&
					expectedResult.BuyerAddress == result.BuyerAddress &&
					expectedResult.ProductID == result.ProductID &&
					expectedResult.ProductSKU == result.ProductSKU &&
					expectedResult.ProductName == result.ProductName &&
					expectedResult.ProductPrice == result.ProductPrice &&
					expectedResult.ProductWeight == result.ProductWeight &&
					expectedResult.ProductDescription == result.ProductDescription &&
					expectedResult.ProductStock == result.ProductStock &&
					expectedResult.ProductUserID == result.ProductUserID &&
					expectedResult.ProductUserFullName == result.ProductUserFullName {

					if len(expectedResult.ProductImagesPath) ==
						len(result.ProductImagesPath) {
						totalImageExist := 0
						for _, expectedProductImagePath := range expectedResult.ProductImagesPath {
							imageExist := false
							for _, resultProductImagePath := range result.ProductImagesPath {
								if expectedProductImagePath == resultProductImagePath {
									imageExist = true
								}
							}
							if imageExist {
								totalImageExist++
							}
						}
						if len(expectedResult.ProductImagesPath) ==
							totalImageExist {
							resultExist = true
						}
					}

				}
			}

			if !resultExist {
				t.Errorf("[%s] Expected order with order number %s not found "+
					"or found but there's invalid value", test.TestName,
					expectedResult.OrderNumber)
			}
		}
	}

	// remove all data order after test
	_, err = collections["orders"].DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("There's an error when truncating "+
			"orders collection after test => %s", err)
	}
}

// getTestingCollections get map of testing mongodb collection
func getTestingCollections(ctx context.Context) (map[string]*mongo.Collection, error) {
	collections := make(map[string]*mongo.Collection)

	// get client mongodb connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DBURI))
	if err != nil {
		return nil, err
	}

	// get list database names in mongodb client connection
	DBNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// check if testing database name is exist in list database names
	DBExist := false
	for _, name := range DBNames {
		if name == config.DBNameForModelTest {
			DBExist = true
			break
		}
	}

	if !DBExist {
		return nil, fmt.Errorf("testing database '%s' not exist",
			config.DBNameForModelTest)
	}

	// get database connection
	DB := client.Database(config.DBNameForModelTest)

	// get list collection name
	collectionNames, err := DB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// check if "orders" collection is exist in database
	collectionExist := false
	for _, name := range collectionNames {
		if name == "orders" {
			collectionExist = true
			break
		}
	}

	if !collectionExist {
		return nil, fmt.Errorf("collection 'orders' not exist")
	}

	// put collection orders to map of collection
	collections["orders"] = DB.Collection("orders")

	return collections, nil
}
