# Dockerfile
FROM golang:1.22.5-alpine

# Set the Current Working Directory inside the container
WORKDIR /home/backend/transfer

# Copy the go.mod and go.sum files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download


# Install make
RUN apt-get update && apt-get install -y make

# Install swag
RUN go get -u github.com/swaggo/swag/cmd/swag

# Copy the application code
COPY . .

# Build the application
RUN go build -o main main.go

# Expose the port
EXPOSE 8080

# Run the command to start the application
CMD ["./main"]