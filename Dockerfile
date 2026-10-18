
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/booking ./cmd/booking

FROM alpine:3.22 AS api

WORKDIR /app
COPY --from=builder /out/api ./api
EXPOSE 8080
ENTRYPOINT ["./api"]

FROM alpine:3.22 AS booking

WORKDIR /app
COPY --from=builder /out/booking ./booking
EXPOSE 50051
ENTRYPOINT ["./booking"]
