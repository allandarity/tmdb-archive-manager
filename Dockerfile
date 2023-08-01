FROM golang:1.20.5-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Golang source code into the container
COPY . .

# Build the Golang application
RUN go build -o app ./src

# Set the command to run your Golang application
CMD ["./app"]
