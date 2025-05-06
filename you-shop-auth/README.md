# Auth Service

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
DB_NAME=products
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8080

LOG_LEVEL=1
LOG_FILE_PATH=auth.log

JWT_SECRET_KEY=SUPPER_SECRET_KEY
JWT_DEFAULT_ACCESS_EXPIRE_TIME=5
JWT_DEFAULT_REFRESH_EXPIRE_TIME=10
```
