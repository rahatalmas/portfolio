# Multi-stage build: frontend (Node) → backend (Go) → small runtime image.

FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/index.html frontend/vite.config.ts frontend/tsconfig.json ./
COPY frontend/src ./src
COPY frontend/public ./public
RUN npm run build

FROM golang:1.24-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=backend /server ./server
COPY --from=frontend /app/dist ./dist
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app
ENV ADDR=:8080 FRONTEND_DIR=/app/dist UPLOAD_DIR=/app/uploads
EXPOSE 8080
ENTRYPOINT ["./server"]
