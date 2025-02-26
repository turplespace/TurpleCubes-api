# Build stage
FROM golang:1.23.3-alpine AS build_base
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the application
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o ./turplecubes ./cmd/main.go

# Copy web folder
COPY web /app/web

# Final stage
FROM docker:dind

# Copy built binary outside bin (in /app)
COPY --from=build_base /app/turplecubes /app/turplecubes

# Copy web folder to /app
COPY --from=build_base /app/web/dist /app/turplecubes_web

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /app

# Set entrypoint
CMD ["/entrypoint.sh"]
