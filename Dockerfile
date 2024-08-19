FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/api


FROM alpine:3.18


WORKDIR /root/

COPY --from=builder /go/bin/api .

COPY --from=builder /app/config.toml .
COPY --from=builder /app/payments.toml .

EXPOSE 8080

CMD ["./api"]
