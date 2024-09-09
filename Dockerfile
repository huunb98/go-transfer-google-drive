# Dockerfile
FROM golang:1.22.5-alpine

# Set the Current Working Directory inside the container
WORKDIR /home/backend/transfer

# Copy the go.mod and go.sum files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download


# Install dependencies
RUN apk add --no-cache make git

# Install swag
RUN go get -u github.com/swaggo/swag/cmd/swag

# Install swag and go-migrate
RUN go install github.com/swaggo/swag/cmd/swag@latest \ 
    && go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy the application code
COPY . .

# Build the application binary
RUN  go build -o main main.go

# Run migrations before starting the application
ENTRYPOINT ["sh", "-c", "make migration_up && go run main.go"]

# Expose the port that the app will listen on
EXPOSE 8008