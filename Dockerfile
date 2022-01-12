# Environment to build manager binary
FROM golang:1.17.6-alpine3.15 as build
WORKDIR /porter

RUN apk update && apk add gcc musl-dev

COPY go.mod go.sum ./
COPY /cmd ./cmd

RUN go mod download

RUN go build -ldflags '-w -s' -a -o ./bin/manager ./cmd/manager

# Deployment environment
# ----------------------
FROM alpine:3.15
WORKDIR /porter

RUN apk update && apk add git

COPY --from=build /porter/bin/manager /usr/bin/

ENTRYPOINT [ "manager" ]
