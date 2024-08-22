FROM golang:1.23.0 AS builder

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go build -o main cmd/cors/main.go

FROM golang:1.23.0

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
