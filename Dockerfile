FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/build/main /app/main
COPY ./data/data.csv /app/data/data.csv

CMD ["/app/main"]
