# build stage
FROM golang:1.16.0 AS builder

WORKDIR /app
COPY go.mod /
COPY go.sum /

RUN go mod download

COPY . .

WORKDIR /app/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache zsh

COPY --from=builder /app/cmd/app/main /app/main

WORKDIR /app

ENTRYPOINT ["./main"]

