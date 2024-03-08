FROM golang:1.22 AS build-stage
WORKDIR /app

COPY . .
# Download Go modules
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /adapter-main

# Deploy the application binary into a lean image
FROM alpine:3.19 AS build-release-stage

WORKDIR /

COPY --from=build-stage /adapter-main /adapter-main
COPY translator-plugins.yaml /
ADD plugins /plugins

# Expose the port
EXPOSE 8080

CMD ["/adapter-main"]
