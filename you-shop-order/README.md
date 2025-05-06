# Order Service

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
DB_PORT=5423
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=order_dev
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8083

LOG_LEVEL=1
LOG_FILE_PATH=order.log

AUTH_SERVER_URL=http://localhost:5000/api/v1/auth

MSG_BROKER_HOST=localhost
MSG_BROKER_PORT=5672
MSG_BROKER_USERNAME=guest
MSG_BROKER_PASSWORD=guest
MSG_BROKER_VHOST=/you_shop

RPC_SERVER_CHECK_PRODUCT_QUANTITY=check_product_quantity
RPC_SERVER_GET_TOTAL_PRICE=get_total_price
RPC_SERVER_GET_PRODUCT_BY_IDS=get_product_by_ids
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
