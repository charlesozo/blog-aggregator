# Use an official Golang runtime as a parent image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download any necessary dependencies
RUN go mod download

# Build the Go application
RUN go build -o myapp

# Expose the port that the application will run on
EXPOSE 8080

# Define environment variables for the PostgreSQL connection
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=12345
ENV DB_NAME=rssagg

# Run the Go application
CMD ["./myapp"]
