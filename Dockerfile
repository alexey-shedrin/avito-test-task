FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/configs /root/configs
COPY --from=builder /app/migrations /root/migrations

EXPOSE 8080

CMD [ "./app" ]