# Stage 1: Builder - Compile the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the binary with optimizations:
# - CGO_ENABLED=0: static binary, no C deps
RUN CGO_ENABLED=0 go build -o /app/bin/app ./cmd/review-assigner


# Stage 2: Runtime - Minimal Alpine image for production
FROM alpine:3.22

RUN addgroup -g 1001 -S appgroup \
	&& adduser -S appuser -u 1001 -G appgroup \
	&& mkdir -p /app \
	&& chown -R appuser:appgroup /app

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/bin/app .

USER appuser:appgroup

ENTRYPOINT ["/bin/sh", "-c", "/app/app"]
