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
RUN npx tailwindcss -i ./input.css -o ./website/public/css/style.css --minify

# Go build
FROM --platform=$BUILDPLATFORM golang:alpine AS build-stage
ARG TARGETOS
ARG TARGETARCH
RUN apk add --no-cache --update go gcc g++
RUN if [ "${TARGETARCH}" = "arm64" ]; then \
  wget -P ~ https://musl.cc/aarch64-linux-musl-cross.tgz \
  && tar -xvzf ~/aarch64-linux-musl-cross.tgz -C /usr/local \
  && rm ~/aarch64-linux-musl-cross.tgz; \
  fi
COPY --from=tailwind-stage /tailwind /app
WORKDIR /app
RUN if [ "${TARGETARCH}" = "arm64" ]; then \
  export CC=/usr/local/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc; \
  fi && \
  CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/main -buildvcs=false ./cmd

# Go test
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy
FROM alpine:latest AS deploy-stage
WORKDIR /app
RUN addgroup --system --gid 5000 app && adduser --system --no-create-home --uid 5000 app --ingroup app
COPY --chown=app:app --from=build-stage /app/website/public /app/public
COPY --chown=app:app --from=build-stage /app/main /app/cmd/main
RUN ls -al
USER app
ENV CLONE_IN_MEMORY=true
EXPOSE 5000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10m \
  CMD wget --no-verbose --tries=1 --spider http://localhost:5000/ || exit 1
RUN ls -al
CMD ["./cmd/main"]
