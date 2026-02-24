# Build stage
FROM golang:1.25.4-alpine3.22 AS builder

WORKDIR /app

RUN apk --no-cache add git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/app

# Runtime stage
FROM alpine:3.23

RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bin/app /bin/app

EXPOSE 8080
ENTRYPOINT ["/bin/app"]
