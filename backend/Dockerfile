FROM golang:1.22-alpine AS builder

ARG APP_VERSION="unknown"

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -v -o /app/todolist-app .

# Final image
FROM alpine:latest

ARG APP_VERSION="unknown"
LABEL version="${APP_VERSION}"

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/todolist-app ./
COPY --from=builder /app/migrations ./migrations/

EXPOSE 8080

CMD ["./todolist-app"]

