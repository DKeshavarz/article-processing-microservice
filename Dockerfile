# Build stage
FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o article-server main.go

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/article-server .
COPY proto/ ./proto/
EXPOSE 50051
CMD ["./article-server"] 