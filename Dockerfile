FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app

FROM alpine:latest

WORKDIR /root/

ENV CONFIG_PATH=configs/dev.yml

COPY --from=builder /app/main .
COPY ./configs/dev.yml ./configs/dev.yml

CMD ["./main","--config=configs/dev.yml"]