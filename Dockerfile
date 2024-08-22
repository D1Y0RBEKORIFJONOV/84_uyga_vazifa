FROM golang:1.23.0 AS builder

RUN mkdir app
COPY . /app

WORKDIR /app

RUN go build -o main cmd/app/main.go

FROM golang:1.23.0

WORKDIR /app

COPY --from=builder /app .

CMD ["/app/main"]

