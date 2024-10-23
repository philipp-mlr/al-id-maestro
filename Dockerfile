# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

# Templ generate
FROM ghcr.io/a-h/templ:latest AS templ-stage
COPY --chown=65532:65532 . /templ
WORKDIR /templ
RUN ["templ", "generate"]

# Tailwind build
FROM node:alpine AS tailwind-stage
WORKDIR /tailwind
COPY package*.json /tailwind
RUN npm install
COPY --from=templ-stage /templ /tailwind
RUN npx tailwindcss -i ./input.css -o ./public/css/style.css --minify

# Go build
FROM golang:latest AS build-stage
COPY --from=tailwind-stage /tailwind /app
WORKDIR /app
RUN GOOS=linux go build -o /app/main -buildvcs=false

# Go test
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy
FROM alpine:latest AS deploy-stage
WORKDIR /app
RUN addgroup --system --gid 5000 app && adduser --system --no-create-home --uid 5000 app --ingroup app
COPY --chown=app:app --from=build-stage /app/public /app/public
COPY --chown=app:app --from=build-stage /app/main /app/main
USER app
EXPOSE 8080
HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1
CMD ["./main"]
