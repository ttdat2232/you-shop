# Price Service

## Prerequisites

-   go version > 1.23.4

## How to run

1. Clone the repository
2. Run `go mod tidy` to download the dependencies
3. Run `go run -tags dev ./api/main.go` to start the API server for development
4. Run `go run -tags prod ./api/main.go` to start the API server for production

## Environment

.env file sample:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=prices
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8080

LOG_LEVEL=1
LOG_FILE_PATH=price.log

MSG_BROKER_HOST=localhost
MSG_BROKER_PORT=5672
MSG_BROKER_USERNAME=guest
MSG_BROKER_PASSWORD=guest
MSG_BROKER_VHOST=/you_shop
```
