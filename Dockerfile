# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12
WORKDIR /app
ENV PORT=8080
ENV DB_PATH=/app/data/packs.db
ENV WEB_DIR=/app/web
COPY --from=builder /app/server /app/server
COPY --from=builder /app/web /app/web
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
