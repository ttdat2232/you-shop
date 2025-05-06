# Build stage
FROM golang:1.24rc1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/main -tags prod ./api/main.go

# Final stage
FROM alpine:3.21

WORKDIR /app/bin

COPY --from=builder /app/bin/main .

CMD ["./main"]
