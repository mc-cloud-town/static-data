FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o start

FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build /app/start server

EXPOSE 8080

CMD ["./server"]
