FROM golang:1.23.6-alpine AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o subscriptions_service .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/subscriptions_service .
COPY config.env .

CMD ["./subscriptions_service"]