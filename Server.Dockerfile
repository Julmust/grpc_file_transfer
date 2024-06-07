FROM golang:1.22

WORKDIR /app

# Install go packages
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy files
COPY server/server.go ./
COPY proto ./proto

# Create environment
RUN mkdir files
COPY server/files/test.png ./files

# Build executable
RUN go build -o ft_server

# Run server
ENTRYPOINT ["/app/ft_server"]