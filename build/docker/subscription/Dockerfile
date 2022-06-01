FROM golang:1.14.3-alpine as builder

# Allows to filter & cleanup intermediate build images later
LABEL stage=builder
LABEL maintainer="Akshit Dhar"

ENV GO111MODULE=on

RUN mkdir /isildur

WORKDIR /isildur

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Export port 8080 for the server binary
EXPOSE 8080

RUN GOOS=linux GOARCH=amd64 go build -o /isildur/exec ./cmd/isildur
CMD ["/isildur/exec"]