# You Shop Payment Service

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
DB_NAME=you_shop_payment
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8080

AUTH_SERVER_URL=http://auth-service.com/api/v1/auth

MSG_BROKER_HOST=localhost
MSG_BROKER_PORT=5672
MSG_BROKER_USERNAME=guest
MSG_BROKER_PASSWORD=guest
MSG_BROKER_VHOST=/you_shop

```

### LOG_LEVEL

-   INFO = 1
-   WARN = 2
-   ERROR = 3
-   FATAL = 4
-   PANIC = 5
-   NO_LEVEL = 6
-   DISABLE = 7
-   TRACE = -1
