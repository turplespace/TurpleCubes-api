FROM golang:1.23.3-alpine AS build_base
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download


# Use make to build and run the application
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/turplecubes ./cmd/main.go


FROM docker:dind
COPY --from=build_base /app/bin /app/bin
CMD [ "./app/bin/turplecubes" ]