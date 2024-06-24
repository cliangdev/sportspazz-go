FROM golang:1.22 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/sportspazz cmd/main.go

FROM alpine:latest as run

WORKDIR /app
COPY --from=builder /app/db/migrations /app/db/migrations
COPY --from=builder /bin/sportspazz /app/bin/sportspazz

ENV PORT 8080
EXPOSE $PORT

CMD ["./bin/sportspazz"]
