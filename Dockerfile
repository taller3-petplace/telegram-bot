FROM golang:latest

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
COPY .env .
COPY ./internal/ ./internal
COPY ./src/ ./src

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/main.go

ENTRYPOINT ["/bin/sh"]