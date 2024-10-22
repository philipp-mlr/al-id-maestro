FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod verify

RUN go build -o ./build/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/main .

EXPOSE 8080

CMD ["./main"]
