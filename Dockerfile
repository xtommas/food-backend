# ==================================================================================== #
# BUILD STAGE
# ==================================================================================== #
FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_TIME
ARG VERSION
RUN go build \
    -ldflags="-s -X main.buildTime=${BUILD_TIME} -X main.version=${VERSION}" \
    -o ./bin/api \
    ./cmd/api

# ==================================================================================== #
# FINAL STAGE
# ==================================================================================== #
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/api ./bin/api

EXPOSE 4000

ENTRYPOINT ["./bin/api"]
