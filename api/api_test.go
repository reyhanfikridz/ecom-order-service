/*
Package api containing API initialization and API route handler
*/
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/reyhanfikridz/ecom-order-service/internal/config"
	"github.com/reyhanfikridz/ecom-order-service/internal/middleware"
	"github.com/reyhanfikridz/ecom-order-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestMain do some test before and after all testing in the package
func TestMain(m *testing.M) {
	ctx := context.Background()

	// init all config before can be used
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("There's an error when initialize config => %s", err)
	}

	// get testing app
	u := middleware.User{}
	a, err := GetTestingAPI(u)
	if err != nil {
		log.Fatalf("There's an error when initialize "+
			"testing API => %s", err)
	}

	// remove all data in all collection before run all test
	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections before run all test => %s", err)
		}
	}

	// run all testing
	m.Run()

	// remove all data in all collection before run all test
	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// TestInitCollections test InitCollections
func TestInitCollections(t *testing.T) {
	a := API{}

	err := a.InitCollections(config.DBName)
	if err != nil {
		t.Errorf("Expected initialize connection to database collections success,"+
			" but connection failed => %s", err)
	}
}

// TestGetOrdersHandler test GetOrdersHandler
func TestGetOrdersHandler(t *testing.T) {
	ctx := context.Background()

	// get testing API
	u := middleware.User{}
	a, err := GetTestingAPI(u)
	if err != nil {
		t.Errorf("There's an error when getting testing API => %s",
			err)
	}

	// create orders
	orders := []model.Order{
		{
			Status:             "in-cart",
			Qty:                2,
			TotalPrice:         2000001,
			BuyerID:            1,
			BuyerFullName:      "George Marcus",
			BuyerAddress:       "Buyer Street",
			ProductID:          1,
			ProductSKU:         "testsku",
			ProductName:        "product name",
			ProductPrice:       1000000.50,
			ProductWeight:      1.5,
			ProductDescription: "product description",
			ProductStock:       100,
			ProductUserID:      10,
			ProductImagesPath:  []string{"product 1.1.jpg", "product 1.2.jpg"},
		},
		{
			Status:             "done",
			Qty:                2,
			TotalPrice:         4000001,
			BuyerID:            1,
			BuyerFullName:      "George Marcus",
			BuyerAddress:       "Buyer Street",
			ProductID:          2,
			ProductSKU:         "testsku2",
			ProductName:        "product name 2",
			ProductPrice:       2000000.50,
			ProductWeight:      2.5,
			ProductDescription: "product description 2",
			ProductStock:       200,
			ProductUserID:      20,
			ProductImagesPath:  []string{"product 2.1.jpg", "product 2.2.jpg", "product 2.3.jpg"},
		},
		{
			Status:             "in-cart",
			Qty:                2,
			TotalPrice:         6000001,
			BuyerID:            2,
			BuyerFullName:      "Linda",
			BuyerAddress:       "Buyer Street 2",
			ProductID:          3,
			ProductSKU:         "testsku 3",
			ProductName:        "product name 3",
			ProductPrice:       3000000.50,
			ProductWeight:      3.5,
			ProductDescription: "product description 3",
			ProductStock:       300,
			ProductUserID:      30,
			ProductImagesPath:  []string{},
		},
	}

	// loop orders
	for i := range orders {
		// insert order into database collection
		orders[i], err = model.InsertOrder(ctx, a.Collections["orders"],
			orders[i])
		if err != nil {
			t.Fatalf("There's an error when insert order => %s",
				err)
		}
	}

	// create testing table
	testTable := []struct {
		TestName        string
		Filter          map[string]string
		User            middleware.User
		ExpectedStatus  int
		ExpectedResults []model.Order
	}{
		{
			TestName:        "Test Get All Order",
			Filter:          nil,
			User:            middleware.User{Role: "buyer"},
			ExpectedStatus:  http.StatusOK,
			ExpectedResults: orders,
		},
		{
			TestName:        "Test Get All Order By Buyer ID <1>",
			Filter:          map[string]string{"buyer_id": "1"},
			User:            middleware.User{Role: "buyer"},
			ExpectedStatus:  http.StatusOK,
			ExpectedResults: []model.Order{orders[0], orders[1]},
		},
		{
			TestName:        "Test Get All Order By Buyer ID <1> and Status <in-cart>",
			Filter:          map[string]string{"buyer_id": "1", "status": "in-cart"},
			User:            middleware.User{Role: "buyer"},
			ExpectedStatus:  http.StatusOK,
			ExpectedResults: []model.Order{orders[0]},
		},
		{
			TestName:        "Test Get All Order By Product User ID <30>",
			Filter:          map[string]string{"product_user_id": "30"},
			User:            middleware.User{Role: "buyer"},
			ExpectedStatus:  http.StatusOK,
			ExpectedResults: []model.Order{orders[2]},
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// get testing API for get products
		a, err = GetTestingAPI(test.User)
		if err != nil {
			t.Errorf("There's an error when getting testing API => %s",
				err)
		}

		// get url params
		params := url.Values{}
		for key, value := range test.Filter {
			params.Add(key, value)
		}

		// create and run request
		req := httptest.NewRequest("GET", "/", nil)
		req.URL.RawQuery = params.Encode()
		if err != nil {
			t.Errorf("[%s] There's an error when creating "+
				"request API get orders => %s",
				test.TestName, err)
		}

		response := httptest.NewRecorder()
		echoCtx := a.Echo.NewContext(req, response)
		echoCtx.Set("user", test.User)
		err = a.GetOrdersHandler(echoCtx)
		if err != nil {
			t.Errorf("[%s] Expected API call success, but got error => %s",
				test.TestName, err)
		}

		// check response status
		if response.Code != test.ExpectedStatus {
			t.Errorf("[%s] Expected status %d got %d",
				test.TestName, test.ExpectedStatus, response.Code)

			var resp map[string]string
			err = json.NewDecoder(response.Body).Decode(&resp)
			if err != nil {
				t.Errorf("[%s] There's an error when unmarshal body response => %s",
					test.TestName, err)
			}

			t.Error(resp)

		} else if response.Code == test.ExpectedStatus &&
			response.Code != http.StatusForbidden {
			// get response data (product)
			var results []model.Order
			err = json.NewDecoder(response.Body).Decode(&results)
			if err != nil {
				t.Errorf("[%s] There's an error when unmarshal body response => %s",
					test.TestName, err)
			}

			// check response data length
			if len(test.ExpectedResults) != len(results) {
				t.Errorf("[%s] Expected length data %d, but got %d",
					test.TestName, len(test.ExpectedResults),
					len(results))
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
						expectedResult.ProductUserID == result.ProductUserID {

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
	}

	// remove all data in all collection after test
	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// TestAddOrderHandler test AddOrderHandler
func TestAddOrderHandler(t *testing.T) {
	// initialize testing table
	testTable := []struct {
		TestName       string
		FormData       map[string]string
		User           middleware.User
		ExpectedOrder  model.Order
		ExpectedStatus int
	}{
		{
			TestName: "Test Add Order Success",
			FormData: map[string]string{
				"status":              "in-cart",
				"qty":                 "2",
				"total_price":         "2000000.50",
				"buyer_full_name":     "George Marcus",
				"buyer_address":       "Buyer Street",
				"product_name":        "Product 1",
				"product_price":       "1000000.50",
				"product_weight":      "1.5",
				"product_description": "Product description",
				"product_stock":       "100",
			},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedOrder: model.Order{
				Status:             "in-cart",
				Qty:                2,
				BuyerID:            1,
				BuyerFullName:      "George Marcus",
				BuyerAddress:       "Buyer Street",
				TotalPrice:         2000000.50,
				ProductName:        "Product 1",
				ProductPrice:       1000000.50,
				ProductWeight:      1.5,
				ProductDescription: "Product description",
				ProductStock:       100,
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			TestName: "Test Add Order Forbidden",
			FormData: map[string]string{
				"status":              "in-cart",
				"qty":                 "2",
				"total_price":         "2000000.50",
				"buyer_full_name":     "George Marcus",
				"buyer_address":       "Buyer Street",
				"product_name":        "Product 1",
				"product_price":       "1000000.50",
				"product_weight":      "1.5",
				"product_description": "Product description",
				"product_stock":       "100",
			},
			User: middleware.User{
				ID:   1,
				Role: "seller",
			},
			ExpectedOrder:  model.Order{},
			ExpectedStatus: http.StatusForbidden,
		},
		{
			TestName: "Test Add Order Bad Request",
			FormData: map[string]string{
				"status":              "",
				"qty":                 "2",
				"total_price":         "2000000.50",
				"buyer_full_name":     "George Marcus",
				"buyer_address":       "Buyer Street",
				"product_name":        "Product 1",
				"product_price":       "1000000.50",
				"product_weight":      "1.5",
				"product_description": "Product description",
			},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedOrder:  model.Order{},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// initialize testing API
		a, err := GetTestingAPI(test.User)
		if err != nil {
			t.Errorf("[%s] There's an error when getting testing API => %s",
				test.TestName, err.Error())
		}

		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("[%s] There's an error when creating "+
					"bytes buffer form data => %s",
					test.TestName, err.Error())
			}

			_, err = io.Copy(fw, strings.NewReader(r))
			if err != nil {
				t.Errorf("[%s] There's an error when creating "+
					"bytes buffer form data => %s",
					test.TestName, err.Error())
			}
		}
		w.Close()

		// create and run request
		req := httptest.NewRequest("POST", "/", &bFormData)
		req.Header.Set("Content-Type", w.FormDataContentType())
		if err != nil {
			t.Errorf("[%s] There's an error when creating "+
				"request API add orders => %s",
				test.TestName, err)
		}

		response := httptest.NewRecorder()
		echoCtx := a.Echo.NewContext(req, response)
		echoCtx.Set("user", test.User)
		err = a.AddOrderHandler(echoCtx)
		if err != nil {
			t.Errorf("[%s] Expected API call success, but got error => %s",
				test.TestName, err)
		}

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("[%s] Expected status %d got %d",
				test.TestName, test.ExpectedStatus, response.Code)

			var resp map[string]string
			err = json.NewDecoder(response.Body).Decode(&resp)
			if err != nil {
				t.Errorf("[%s] There's an error when unmarshal body response => %s",
					test.TestName, err.Error())
			}

			t.Error(resp)
		} else {
			if test.ExpectedStatus == http.StatusCreated {
				// get response data (product)
				var respPInfo model.Order
				err = json.NewDecoder(response.Body).Decode(&respPInfo)
				if err != nil {
					t.Errorf("[%s] There's an error when unmarshal body response => %s",
						test.TestName, err.Error())
				}

				// check result
				if respPInfo.ID == primitive.NilObjectID {
					t.Errorf("[%s] Expected id not nil, but got nil", test.TestName)
				}
				if test.ExpectedOrder.ProductName != respPInfo.ProductName {
					t.Errorf("[%s] Expected ProductName '%s', but got ProductName '%s'",
						test.TestName, test.ExpectedOrder.ProductName,
						respPInfo.ProductName)
				}
				if test.ExpectedOrder.ProductPrice != respPInfo.ProductPrice {
					t.Errorf("[%s] Expected ProductPrice %f, but got ProductPrice %f",
						test.TestName, test.ExpectedOrder.ProductPrice,
						respPInfo.ProductPrice)
				}
				if test.ExpectedOrder.ProductWeight != respPInfo.ProductWeight {
					t.Errorf("[%s] Expected ProductWeight %f, but got ProductWeight %f",
						test.TestName, test.ExpectedOrder.ProductWeight,
						respPInfo.ProductWeight)
				}
				if test.ExpectedOrder.ProductDescription != respPInfo.ProductDescription {
					t.Errorf("[%s] Expected ProductDescription '%s', "+
						"but got ProductDescription '%s'",
						test.TestName, test.ExpectedOrder.ProductDescription,
						respPInfo.ProductDescription)
				}
				if test.ExpectedOrder.ProductStock != respPInfo.ProductStock {
					t.Errorf("[%s] Expected ProductStock %d, but got ProductStock %d",
						test.TestName, test.ExpectedOrder.ProductStock,
						respPInfo.ProductStock)
				}
				if test.ExpectedOrder.BuyerID != respPInfo.BuyerID {
					t.Errorf("[%s] Expected BuyerID %d, but got BuyerID %d",
						test.TestName, test.ExpectedOrder.BuyerID,
						respPInfo.BuyerID)
				}
				if test.ExpectedOrder.BuyerFullName != respPInfo.BuyerFullName {
					t.Errorf("[%s] Expected BuyerFullName %s, but got BuyerFullName %s",
						test.TestName, test.ExpectedOrder.BuyerFullName,
						respPInfo.BuyerFullName)
				}
				if test.ExpectedOrder.BuyerAddress != respPInfo.BuyerAddress {
					t.Errorf("[%s] Expected BuyerAddress %s, but got BuyerAddress %s",
						test.TestName, test.ExpectedOrder.BuyerAddress,
						respPInfo.BuyerAddress)
				}
			}
		}
	}

	// remove all data in all collection after test
	a, err := GetTestingAPI(middleware.User{})
	if err != nil {
		t.Errorf("There's an error when getting testing API => %s",
			err.Error())
	}

	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(a.Ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// TestUpdateOrderHandler test UpdateOrderHandler
//
// Required for test: model.InsertOrder
func TestUpdateOrderHandler(t *testing.T) {
	// insert testing data
	a, err := GetTestingAPI(middleware.User{
		ID:   1,
		Role: "buyer",
	})
	if err != nil {
		t.Fatalf("There's an error when getting testing API => %s",
			err.Error())
	}
	oCreate := model.Order{
		Status:             "in-cart",
		Qty:                2,
		BuyerID:            1,
		BuyerFullName:      "George Marcus",
		BuyerAddress:       "Buyer Street",
		TotalPrice:         2000000.50,
		ProductName:        "Product 1",
		ProductPrice:       1000000.50,
		ProductWeight:      1.5,
		ProductDescription: "Product description",
		ProductStock:       100,
	}

	oCreate, err = model.InsertOrder(a.Ctx, a.Collections["orders"], oCreate)
	if err != nil {
		t.Fatalf("There's an error when creating testing data for testing "+
			"update data => %s", err.Error())
	}

	// initialize testing table
	testTable := []struct {
		TestName       string
		Filter         map[string]string
		FormData       map[string]string
		User           middleware.User
		ExpectedStatus int
	}{
		{
			TestName: "Test Update Order Success",
			Filter: map[string]string{
				"order_number": oCreate.OrderNumber,
			},
			FormData: map[string]string{
				"status":      "waiting-for-payment",
				"qty":         "3",
				"total_price": "3000001.00",
			},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			TestName: "Test Update Order Bad Request",
			Filter:   map[string]string{},
			FormData: map[string]string{
				"status":      "",
				"qty":         "3",
				"total_price": "3000001.00",
			},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// initialize testing API
		a, err := GetTestingAPI(test.User)
		if err != nil {
			t.Errorf("[%s] There's an error when getting testing API => %s",
				test.TestName, err.Error())
		}

		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("[%s] There's an error when creating "+
					"bytes buffer form data => %s",
					test.TestName, err.Error())
			}

			_, err = io.Copy(fw, strings.NewReader(r))
			if err != nil {
				t.Errorf("[%s] There's an error when creating "+
					"bytes buffer form data => %s",
					test.TestName, err.Error())
			}
		}
		w.Close()

		// get url params
		params := url.Values{}
		for key, value := range test.Filter {
			params.Add(key, value)
		}

		// create and run request
		req := httptest.NewRequest("PUT", "/", &bFormData)
		req.URL.RawQuery = params.Encode()
		req.Header.Set("Content-Type", w.FormDataContentType())
		if err != nil {
			t.Errorf("[%s] There's an error when creating "+
				"request API add orders => %s",
				test.TestName, err)
		}

		response := httptest.NewRecorder()
		echoCtx := a.Echo.NewContext(req, response)
		echoCtx.Set("user", test.User)
		err = a.UpdateOrderHandler(echoCtx)
		if err != nil {
			t.Errorf("[%s] Expected API call success, but got error => %s",
				test.TestName, err)
		}

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("[%s] Expected status %d got %d",
				test.TestName, test.ExpectedStatus, response.Code)

			var resp map[string]string
			err = json.NewDecoder(response.Body).Decode(&resp)
			if err != nil {
				t.Errorf("[%s] There's an error when unmarshal body response => %s",
					test.TestName, err.Error())
			}

			t.Error(resp)
		}
	}

	// remove all data in all collection after test
	a, err = GetTestingAPI(middleware.User{})
	if err != nil {
		t.Errorf("There's an error when getting testing API => %s",
			err.Error())
	}

	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(a.Ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// TestDeleteOrderHandler test DeleteOrderHandler
//
// Required for test: model.InsertOrder
func TestDeleteOrderHandler(t *testing.T) {
	// insert testing data
	a, err := GetTestingAPI(middleware.User{
		ID:   1,
		Role: "buyer",
	})
	if err != nil {
		t.Fatalf("There's an error when getting testing API => %s",
			err.Error())
	}
	oCreate := model.Order{
		Status:             "in-cart",
		Qty:                2,
		BuyerID:            1,
		BuyerFullName:      "George Marcus",
		BuyerAddress:       "Buyer Street",
		TotalPrice:         2000000.50,
		ProductName:        "Product 1",
		ProductPrice:       1000000.50,
		ProductWeight:      1.5,
		ProductDescription: "Product description",
		ProductStock:       100,
	}

	oCreate, err = model.InsertOrder(a.Ctx, a.Collections["orders"], oCreate)
	if err != nil {
		t.Fatalf("There's an error when creating testing data for testing "+
			"update data => %s", err.Error())
	}

	// initialize testing table
	testTable := []struct {
		TestName       string
		Filter         map[string]string
		User           middleware.User
		ExpectedStatus int
	}{
		{
			TestName: "Test Delete Order Success",
			Filter: map[string]string{
				"order_number": oCreate.OrderNumber,
			},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			TestName: "Test Delete Order Bad Request",
			Filter:   map[string]string{},
			User: middleware.User{
				ID:   1,
				Role: "buyer",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// initialize testing API
		a, err := GetTestingAPI(test.User)
		if err != nil {
			t.Errorf("[%s] There's an error when getting testing API => %s",
				test.TestName, err.Error())
		}

		// get url params
		params := url.Values{}
		for key, value := range test.Filter {
			params.Add(key, value)
		}

		// create and run request
		req := httptest.NewRequest("DELETE", "/", nil)
		req.URL.RawQuery = params.Encode()
		if err != nil {
			t.Errorf("[%s] There's an error when creating "+
				"request API add orders => %s",
				test.TestName, err)
		}

		response := httptest.NewRecorder()
		echoCtx := a.Echo.NewContext(req, response)
		echoCtx.Set("user", test.User)
		err = a.DeleteOrderHandler(echoCtx)
		if err != nil {
			t.Errorf("[%s] Expected API call success, but got error => %s",
				test.TestName, err)
		}

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("[%s] Expected status %d got %d",
				test.TestName, test.ExpectedStatus, response.Code)

			var resp map[string]string
			err = json.NewDecoder(response.Body).Decode(&resp)
			if err != nil {
				t.Errorf("[%s] There's an error when unmarshal body response => %s",
					test.TestName, err.Error())
			}

			t.Error(resp)
		}
	}

	// remove all data in all collection after test
	a, err = GetTestingAPI(middleware.User{})
	if err != nil {
		t.Errorf("There's an error when getting testing API => %s",
			err.Error())
	}

	for _, collection := range a.Collections {
		_, err = collection.DeleteMany(a.Ctx, bson.D{})
		if err != nil {
			log.Fatalf("There's an error when truncating "+
				"mongodb collections after run all test => %s", err)
		}
	}
}

// GetTestingAPI get API for testing
func GetTestingAPI(u middleware.User) (API, error) {
	a := API{}
	a.Echo = echo.New()

	// init mongodb database collections
	err := a.InitCollections(config.DBNameForAPITest)
	if err != nil {
		return a, err
	}

	return a, nil
}
