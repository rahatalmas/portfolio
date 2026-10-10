# Backend only: JSON API + admin + uploads. The React site is built and
# hosted separately (any static host) and points at this server via
# VITE_API_BASE; list the site's origin in CORS_ORIGINS here.

FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /server ./server
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app
ENV ADDR=:8080 UPLOAD_DIR=/app/uploads
EXPOSE 8080
ENTRYPOINT ["./server"]
