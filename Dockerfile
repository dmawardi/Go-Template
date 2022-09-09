FROM golang:1.19-bullseye

# Set work dir of Docker container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# RUN go install github.com/dmawardi/Go-Template@latest
# Copy current directory and copy to Docker working directory
COPY . .

# Build the Go app
RUN go build -o ./out/go-template ./cmd

# Expose port 8080
EXPOSE 8080


CMD ["./out/go-template"]