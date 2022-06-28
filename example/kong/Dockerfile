FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o validate-protobuf ./plugin/kong/
RUN /app/validate-protobuf -dump

FROM kong:2.8.1-alpine

COPY --from=builder /app/validate-protobuf /usr/local/bin/validate-protobuf
RUN /usr/local/bin/validate-protobuf -dump
