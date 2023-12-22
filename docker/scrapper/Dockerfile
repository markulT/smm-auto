FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN go mod download

# Build the Go application
RUN go build -o golearn .

# Expose the port the app runs on
EXPOSE 8080

# Define the command to run your application
CMD ["./golearn"]