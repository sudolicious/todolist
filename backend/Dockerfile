FROM golang:1.22-alpine AS builder

ARG APP_VERSION="unknown"

RUN apk update && apk add --no-cache git
ENV GOPATH=/

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -v -o /app/todolist-app .

# Final image
FROM alpine:latest

ARG APP_VERSION="unknown"
#meta-data
LABEL version="${APP_VERSION}"

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/todolist-app ./
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/.env ./

EXPOSE 8080
CMD ["./todolist-app"]
