# Use the official Go image as a base image
FROM golang:1.21.3 AS build

# Set the working directory inside the container
WORKDIR /kyn-project

# Copy the go.mod and go.sum files to the container's WORKDIR
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project to the container's WORKDIR
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

### Start a new stage from scratch ###
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the build stage to the current stage
COPY --from=build /kyn-project/main .

# Command to run the executable
CMD ["./main"]
