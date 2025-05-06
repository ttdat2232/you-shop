# Image Service

## Prerequisites

-   go version > 1.23.4

## How to run

1. Clone the repository
2. Run `go mod tidy` to download the dependencies
3. Run `go run -tags dev ./cmd/main.go` to start the API server for development
4. Run `go run -tags prod ./cmd/main.go` to start the API server for production

## Environment

.env file sample:

```
DB_HOST=localhost
DB_PORT=30432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=images_dev
DB_SSLMODE=disable

API_SERVER_HOST=localhost
API_SERVER_PORT=8081

AUTH_SERVER_URL=http://localhost5000/api/v1/auth/

LOG_LEVEL=-1
LOG_FILE_PATH=image.log

MSG_BROKER_HOST=localhost
MSG_BROKER_PORT=5672
MSG_BROKER_USERNAME=guest
MSG_BROKER_PASSWORD=guest
MSG_BROKER_VHOST=/you_shop

S3_PROXY_USERNAME=admin
S3_PROXY_PASSWORD=password
S3_PROXY_HOST=http://host.example.com
S3_PROXY_PORT=30080
S3_PROXY_FOLDER=image
```
