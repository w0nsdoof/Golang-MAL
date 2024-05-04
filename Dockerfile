# use official Golang image
FROM golang:1.21.5

# set working directory
WORKDIR /myanimelist

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o myanimelist ./cmd/myAnimeList

#EXPOSE the port
EXPOSE 8081

# Run the executable
CMD ["./myanimelist"]