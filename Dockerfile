FROM golang:1.24-alpine AS BUILDER
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server main.go

FROM alpine
WORKDIR /app

COPY --from=BUILDER /app/server .
COPY --from=BUILDER /app/migrations ./migrations

EXPOSE 9000
CMD ["./server"]