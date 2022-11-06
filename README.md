# ecom-order-service

### ECOM summary:
ECOM is a simple E-Commerce website builded with Go backend microservices and Django frontend. Disclaimer! I have zero real experience in building E-Commerce system, so if the system is really bad, I apologized in advance. This is just my personal project using Go microservices. You can use all the code of this project as a template for real E-Commerce in the future if you like it. Disclaimer again! I also not a frontend specialist, so I just use a free template I found in the internet and an original bootstrap template.

### Repository summary:
This is a microservice for ECOM that related to customer orders CRUD.

### Requirements:
1. go (recommended: v1.18.4)
2. mongodb (recommended: v6.0.1)

### Microservice requirements:
1. ecom-account-service (must: https://github.com/reyhanfikridz/ecom-account-service/tree/release-1)
2. ecom-product-service (must: https://github.com/reyhanfikridz/ecom-product-service/tree/release-1)

### Steps to run the server:
1. install all requirements
2. install and run all microservice requirements
3. clone repository at directory `$GOPATH/src/github.com/`
4. install required go library with `go mod download` then `go mod vendor` at repository root directory (same level as README.md)
5. create file .env at repository root directory (same level as README.md) with contents:

```
ECOM_ORDER_SERVICE_DB_URI=<mongodb database uri, example: mongodb://localhost:27017>
ECOM_ORDER_SERVICE_DB_NAME=<database name, example: ecomorderservicedb>
ECOM_ORDER_SERVICE_DB_NAME_FOR_API_TEST=<database name for overall api testing, example: ecomorderserviceapitestdb>
ECOM_ORDER_SERVICE_DB_NAME_FOR_MODEL_TEST=<database name for model crud testing, example: ecomorderservicemodeltestdb>

ECOM_ORDER_SERVICE_URL=<this service url, example: :8030>
ECOM_ORDER_SERVICE_FRONTEND_URL=<ecom frontend url, example: http://127.0.0.1:8000>
ECOM_ORDER_SERVICE_ACCOUNT_SERVICE_URL=<ecom frontend url, example: http://127.0.0.1:8010>
ECOM_ORDER_SERVICE_PRODUCT_SERVICE_URL=<ecom frontend url, example: http://127.0.0.1:8020>
```

6. create mongodb databases with name same as in .env file
7. test server first with `go test ./...` to make sure server works fine
8. run server with `go run ./...`
