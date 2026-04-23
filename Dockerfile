FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/kopesa-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/kopesa-migrate ./cmd/migrate

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/kopesa-api /app/kopesa-api
COPY --from=builder /out/kopesa-migrate /app/kopesa-migrate
COPY migrations /app/migrations
COPY api /app/api

ENV PORT=8080
EXPOSE 8080

CMD ["/app/kopesa-api"]
