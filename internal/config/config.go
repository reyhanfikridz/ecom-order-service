/*
Package config collection of configuration
*/
package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	DBURI              string
	DBName             string
	DBNameForAPITest   string
	DBNameForModelTest string

	FrontendURL       string
	AccountServiceURL string
	ProductServiceURL string
)

// InitConfig initialize all config variable from environment variable
func InitConfig() error {
	// load all values from .env file into the system
	// .env file must be at root directory (same level as go.mod file)
	err := godotenv.Load(os.ExpandEnv(
		"$GOPATH/src/github.com/reyhanfikridz/ecom-order-service/.env"))
	if err != nil {
		return err
	}

	// set all config variable after all environment variable loaded
	DBURI = os.Getenv("ECOM_ORDER_SERVICE_DB_URI")
	DBName = os.Getenv("ECOM_ORDER_SERVICE_DB_NAME")
	DBNameForAPITest = os.Getenv("ECOM_ORDER_SERVICE_DB_NAME_FOR_API_TEST")
	DBNameForModelTest = os.Getenv("ECOM_ORDER_SERVICE_DB_NAME_FOR_MODEL_TEST")

	FrontendURL = os.Getenv("ECOM_ORDER_SERVICE_FRONTEND_URL")
	AccountServiceURL = os.Getenv("ECOM_ORDER_SERVICE_ACCOUNT_SERVICE_URL")
	ProductServiceURL = os.Getenv("ECOM_ORDER_SERVICE_PRODUCT_SERVICE_URL")

	return nil
}
