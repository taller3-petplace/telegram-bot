FROM golang:latest

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
COPY .env .
COPY ./internal/ ./internal
COPY ./src/ ./src

RUN go mod download

ENTRYPOINT ["/bin/sh"]