# You Shop Shipment Service

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
DB_HOST=linuxserver
DB_PORT=30432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shipment_dev
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8085

LOG_LEVEL=-1

AUTH_SERVER_URL=http://localhost:5000/api/v1/auth
UPLOAD_SERVER_URL=http://localhost:8081/api/v1/images

MSG_BROKER_HOST=linuxserver
MSG_BROKER_PORT=30672
MSG_BROKER_USERNAME=youshop
MSG_BROKER_PASSWORD=1qaz!QAZ
MSG_BROKER_VHOST=/you_shop_dev

GHN_BASE_URL=https://online-gateway.ghn.vn/shiip/public-api
GHN_TOKEN=a5964f37-0248-11f0-821b-7e7ee35e0791
GHN_SHOP_ID=5687380
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
