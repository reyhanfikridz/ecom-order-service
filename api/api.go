/*
Package api containing API initialization and API route handler
*/
package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reyhanfikridz/ecom-order-service/internal/config"
	"github.com/reyhanfikridz/ecom-order-service/internal/middleware"
	"github.com/reyhanfikridz/ecom-order-service/internal/model"
	"github.com/reyhanfikridz/ecom-order-service/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API contain context, map of mongodb collection, and echo router
type API struct {
	Ctx         context.Context
	Collections map[string]*mongo.Collection
	Echo        *echo.Echo
}

// InitCollections initialize API connection to mongodb collections
func (a *API) InitCollections(DBName string) error {
	// make collections
	a.Collections = make(map[string]*mongo.Collection)

	// get client mongodb connection
	client, err := mongo.Connect(a.Ctx, options.Client().ApplyURI(config.DBURI))
	if err != nil {
		return err
	}

	// get list database names in mongodb client connection
	DBNames, err := client.ListDatabaseNames(a.Ctx, bson.M{})
	if err != nil {
		return err
	}

	// check if database name is exist in list database names
	DBExist := false
	for _, name := range DBNames {
		if name == DBName {
			DBExist = true
			break
		}
	}

	if !DBExist {
		return fmt.Errorf("database '%s' not exist", DBName)
	}

	// get database connection
	DB := client.Database(DBName)

	// get list collection name
	collectionNames, err := DB.ListCollectionNames(a.Ctx, bson.M{})
	if err != nil {
		return err
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
		return fmt.Errorf("collection 'orders' not exist")
	}

	// put collection orders to map of collection
	a.Collections["orders"] = DB.Collection("orders")

	return nil
}

// InitRouter initialize echo router for API
func (a *API) InitRouter() {
	a.Echo = echo.New()

	// add middleware CORS and logger to all route
	a.Echo.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{
			config.FrontendURL, config.AccountServiceURL, config.ProductServiceURL,
		},
	}))
	a.Echo.Use(echomiddleware.Logger())

	// create main router group (prefix: "/api") with middleware authorization
	mainRouter := a.Echo.Group("/api", middleware.AuthorizationMiddleware)

	//// route add order
	mainRouter.POST("/order/", a.AddOrderHandler)

	//// route get orders
	mainRouter.GET("/orders/", a.GetOrdersHandler)

	//// route update order
	mainRouter.PUT("/order/", a.UpdateOrderHandler)

	//// route delete order
	mainRouter.DELETE("/order/", a.DeleteOrderHandler)
}

// GetOrdersHandler route handler for get orders (Method: GET, User: all)
func (a *API) GetOrdersHandler(c echo.Context) error {
	// get user data
	tmpU := c.Get("user")
	_, ok := tmpU.(middleware.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "user data invalid",
		})
	}

	// get filter
	filter := bson.M{}

	//// get filter buyerID
	buyerID, err := strconv.Atoi(c.QueryParam("buyer_id"))
	if err == nil {
		filter["buyer_id"] = buyerID
	}

	//// get filter status
	if c.QueryParam("status") != "" {
		filter["status"] = c.QueryParam("status")
	}

	//// get filter productUserID
	productUserID, err := strconv.Atoi(c.QueryParam("product_user_id"))
	if err == nil {
		filter["product_user_id"] = productUserID
	}

	// get orders from orders collection
	orders, err := model.GetOrders(a.Ctx, a.Collections["orders"],
		filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf(
				"There's an error when getting the orders data => %s",
				err),
		})
	}

	return c.JSON(http.StatusOK, orders)
}

// AddOrderHandler route handler for add order (Method: POST, User: buyer)
func (a *API) AddOrderHandler(c echo.Context) error {
	// get user data
	tmpU := c.Get("user")
	u, ok := tmpU.(middleware.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "user data invalid",
		})
	}

	// check user role is buyer
	if u.Role != "buyer" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "user doesn't have authority to access this API",
		})
	}

	// set order that need to be inserted to database
	var o model.Order
	err := c.Bind(&o)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("Order data not completed/invalid => %s", err),
		})
	}
	o.BuyerID = u.ID

	// validate order data
	err = validator.IsOrderValid(o)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("Order data not completed/invalid => %s", err),
		})
	}

	// insert order to database
	o, err = model.InsertOrder(a.Ctx, a.Collections["orders"], o)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf(
				"There's an error when inserting order data => %s",
				err),
		})
	}

	return c.JSON(http.StatusCreated, o)
}

// UpdateOrderHandler route handler for update order (Method: PUT, User: all)
func (a *API) UpdateOrderHandler(c echo.Context) error {
	// get user data
	tmpU := c.Get("user")
	_, ok := tmpU.(middleware.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "user data invalid",
		})
	}

	// set order that need to be updated to database
	var o model.Order
	err := c.Bind(&o)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("Order data not completed/invalid => %s", err),
		})
	}

	// set filter (for now only order number)
	filter := bson.M{}
	if c.QueryParam("order_number") == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "order_number empty/not found",
		})
	}
	filter["order_number"] = c.QueryParam("order_number")

	// update order in database
	err = model.UpdateOrder(a.Ctx, a.Collections["orders"], filter, o)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf(
				"There's an error when updating order data => %s",
				err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Update order success!",
	})
}

// DeleteOrderHandler route handler for delete order (Method: DELETE, User: all)
func (a *API) DeleteOrderHandler(c echo.Context) error {
	// get user data
	tmpU := c.Get("user")
	_, ok := tmpU.(middleware.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "user data invalid",
		})
	}

	// set filter (for now only order number)
	filter := bson.M{}
	if c.QueryParam("order_number") == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "order_number empty/not found",
		})
	}
	filter["order_number"] = c.QueryParam("order_number")

	// update order in database
	err := model.DeleteOrder(a.Ctx, a.Collections["orders"], filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf(
				"There's an error when deleting order data => %s",
				err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Delete order success!",
	})
}
