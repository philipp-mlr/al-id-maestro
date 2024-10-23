# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app
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
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main -buildvcs=false

# Go test
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy
FROM alpine:latest AS deploy-stage
WORKDIR /app
RUN addgroup --system --gid 5000 app && adduser --system --no-create-home --uid 5000 app --ingroup app
COPY --chown=app:app --from=build-stage /app/public /app/public
COPY --chown=app:app --from=build-stage /app/main /app/main
COPY --chown=app:app --from=build-stage /app/translation.json /app/translation.json
USER app
ENV EMAIL_FROM=contact-request@philipp-cloud.io
ENV EMAIL_TO=hello@philipp.software
ENV EMAIL_USERNAME=webmaster@philipp-cloud.io
ENV EMAIL_HOST=smtp.hostinger.com
ENV EMAIL_PORT=465
EXPOSE 1323
HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:1323/ || exit 1
CMD ["./main"]
