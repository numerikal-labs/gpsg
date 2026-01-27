# Stage 1: Build
FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o gpsg .

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gpsg .
EXPOSE 8080
ENTRYPOINT ["./gpsg"]
