FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Install required system dependencies
RUN apk add --no-cache wget
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz && \
    tar xvzf migrate.tar.gz

#Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY wait-for.sh .
COPY start.sh .
COPY db/migration ./migration

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
