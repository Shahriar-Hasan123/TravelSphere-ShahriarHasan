# Stage 1 — Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git — required by some Go modules
RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o travelsphere .

# Stage 2 — Runtime
FROM alpine:3.20

# ca-certificates — required for HTTPS calls to external APIs
RUN apk --no-cache add ca-certificates wget

# Run as non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/travelsphere .

# Copy required runtime assets
COPY --from=builder /app/views/   ./views/
COPY --from=builder /app/static/  ./static/
COPY --from=builder /app/conf/    ./conf/

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./travelsphere"]