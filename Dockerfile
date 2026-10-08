# syntax=docker/dockerfile:1
FROM golang:1.26.4-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/app ./cmd/app

FROM alpine:3.23
RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /bin/app /bin/app
ENV GIN_MODE=release
USER app:app
EXPOSE 8080
STOPSIGNAL SIGTERM
ENTRYPOINT ["/bin/app"]
