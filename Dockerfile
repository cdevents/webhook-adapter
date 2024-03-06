FROM golang:1.21
WORKDIR /app
COPY . .
# Build the Go application
RUN go build -o main .
# Expose the port
EXPOSE 8080
CMD ["./main"]
