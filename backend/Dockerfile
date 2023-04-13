# Set the base image to the official Go image
FROM golang:1.20

# Set the working directory to /go/src/app
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o ip-lookup-app .

# Expose port 8080
EXPOSE 8080

# Start the Go application
CMD ["./ip-lookup-app"]
