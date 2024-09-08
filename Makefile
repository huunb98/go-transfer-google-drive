# Load the environment variables from the .env file
include .env
export $(shell sed 's/=.*//' .env)

build:
	swag init
	sed -i '/LeftDelim:/d' docs/docs.go
	sed -i '/RightDelim:/d' docs/docs.go
	go build main.go

run:
	go run main.go

# Migration
migration_generate:
	migrate create -ext sql -dir internal/database/migrations -seq ${name}

migration_up:
	migrate -path ./internal/database/migrations -database "mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE_NAME}" up

migration_down:
	migrate -path ./internal/database/migrations -database "mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE_NAME}" down
