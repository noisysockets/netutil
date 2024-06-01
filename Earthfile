VERSION 0.7
FROM golang:1.22-bookworm
WORKDIR /workspace

tidy:
  LOCALLY
  RUN go mod tidy
  RUN go fmt ./...

lint:
  FROM golangci/golangci-lint:v1.57.2
  WORKDIR /workspace
  COPY . .
  RUN golangci-lint run --timeout 5m ./...

test:
  COPY go.* .
  RUN go mod download
  COPY . .
  HOST demo.example.com 127.0.1.1
  HOST demo 127.0.1.1
  RUN --privileged hostname demo.example.com \
    && go test -coverprofile=coverage.out -v ./...
  SAVE ARTIFACT coverage.out AS LOCAL coverage.out