# First: Builder env
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
# CGO Disabled, Linux target static build
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-sensor

# Second: OP Env 
FROM alpine:latest
WORKDIR /root/
# bring exercise file
COPY --from=builder /app/mock-sensor .
CMD ["./mock-sensor"]
