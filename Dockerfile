FROM golang:1.17 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -installsuffix cgo -o app ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates && apk --no-cache add tzdata
WORKDIR /app
COPY --from=0 /app .
CMD mkdir /var/data
ENTRYPOINT ["/app/app"]