/*
Package main the executeable file
*/
package main

import (
	"log"

	"github.com/reyhanfikridz/ecom-order-service/api"
	"github.com/reyhanfikridz/ecom-order-service/internal/config"
)

// main
func main() {
	// init API
	a, err := InitAPI()
	if err != nil {
		log.Fatal(err)
	}

	// serve server
	log.Fatal(a.Echo.Start(":8030"))
}

// InitAPI initialize API
func InitAPI() (api.API, error) {
	a := api.API{}

	// init all config before can be used
	err := config.InitConfig()
	if err != nil {
		return a, err
	}

	// init database
	err = a.InitCollections(config.DBName)
	if err != nil {
		return a, err
	}

	// init router
	a.InitRouter()

	return a, nil
}
