# ecom-order-service

### ECOM summary:
ECOM is a simple E-Commerce website builded with Go backend microservices and Django frontend. Disclaimer! I have zero real experience in building E-Commerce system, so if the system is really bad, I apologized in advance. This is just my personal project using Go microservices. You can use all the code of this project as a template for real E-Commerce in the future if you like it. Disclaimer again! I also not a frontend specialist, so I just use a free template I found in the internet and an original bootstrap template.

### Repository summary:
This is a microservice for ECOM that related to customer orders CRUD.

### Requirements:
1. go (tested: v1.18.4, v1.19.3)
2. mongodb (tested: v6.0.1, v6.0.2)

### Microservice requirements:
1. ecom-account-service (must: https://github.com/reyhanfikridz/ecom-account-service/tree/release-1)
2. ecom-product-service (must: https://github.com/reyhanfikridz/ecom-product-service/tree/release-1)

### Steps to run the server:
1. install all requirements
2. install and run all microservice requirements
3. clone repository with `git clone https://github.com/reyhanfikridz/ecom-order-service` at directory `$GOPATH/src/github.com/reyhanfikridz/`
4. change branch to release-1 with `git checkout release-1` then `git pull origin release-1` at repository root directory (same level as README.md)
5. install required go library with `go mod download` then `go mod vendor` at repository root directory (same level as README.md)
6. create file .env at repository root directory (same level as README.md) with contents:

```
ECOM_ORDER_SERVICE_DB_URI=<mongodb database uri, example: mongodb://localhost:27017>
ECOM_ORDER_SERVICE_DB_NAME=<database name, example: ecomorderservicedb>
ECOM_ORDER_SERVICE_DB_NAME_FOR_API_TEST=<database name for overall api testing, example: ecomorderserviceapitestdb>
ECOM_ORDER_SERVICE_DB_NAME_FOR_MODEL_TEST=<database name for model crud testing, example: ecomorderservicemodeltestdb>

ECOM_ORDER_SERVICE_URL=<this service url, example: :8030>
ECOM_ORDER_SERVICE_FRONTEND_URL=<ecom frontend url, example: http://127.0.0.1:8000>
ECOM_ORDER_SERVICE_ACCOUNT_SERVICE_URL=<ecom account service url, example: http://127.0.0.1:8010>
ECOM_ORDER_SERVICE_PRODUCT_SERVICE_URL=<ecom product service url, example: http://127.0.0.1:8020>
```

7. create mongodb databases with name same as in .env file
8. test server first with `go test ./...` to make sure server works fine
9. run server with `go run ./...`
