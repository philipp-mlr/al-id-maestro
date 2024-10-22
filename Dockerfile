FROM golang:latest

WORKDIR /app

COPY ./handler ./handler
COPY ./model ./model
COPY ./component ./component
COPY ./service ./service
COPY ./public ./public
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

RUN go build -o ./build/main .

EXPOSE 8080

CMD ["./build/main"]
