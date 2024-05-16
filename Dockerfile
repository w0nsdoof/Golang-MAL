# use official Golang image
FROM golang:1.21.5

RUN go install github.com/gobuffalo/pop/soda@latest
RUN apt-get update && apt-get install -y postgresql-client

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