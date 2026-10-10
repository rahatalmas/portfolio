# One image for the whole portfolio: the React site, the JSON API, the admin
# and uploads, all served by the Go binary on :8080. No database: content
# lives in /app/data (mount it as a volume), uploads in /app/uploads.
#
# Stage 1 builds the frontend with Node, stage 2 builds the Go server, and the
# final stage copies both into a small Alpine image.

FROM node:22-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
# The site embeds content.json as its offline fallback, so the build needs it
# at the same relative path as in the repo.
COPY backend/data/content.json /src/backend/data/content.json
COPY frontend/index.html frontend/vite.config.ts frontend/tsconfig*.json ./
COPY frontend/src ./src
COPY frontend/public ./public
RUN npm run build

FROM golang:1.26-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=backend /server ./server
COPY --from=frontend /src/frontend/dist ./dist
# Starter content. A named volume mounted at /app/data is seeded from this
# directory the first time it is created; after that the volume is the truth.
COPY backend/data/content.json ./data/content.json
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app
ENV ADDR=:8080 FRONTEND_DIR=/app/dist DATA_DIR=/app/data UPLOAD_DIR=/app/uploads
EXPOSE 8080
ENTRYPOINT ["./server"]
