FROM golang:1.22-alpine AS build
LABEL authors="waleed"

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

COPY ./ ./

RUN go mod download

RUN go build -v -o go_app ./main.go

RUN chmod +x go_app

FROM alpine:3.17
WORKDIR /app

COPY --from=build /app/go_app ./
COPY .env ./
COPY localhost.pem ./
COPY localhost-key.pem ./


CMD ["./go_app"]
