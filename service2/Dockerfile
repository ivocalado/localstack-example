# Stage 1: Build the Golang application
FROM golang:1.21.6 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Stage 2: Copy the binary to a distroless image
FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
