FROM golang:1.16-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/app .

# Build a small image
FROM alpine:latest

ENV FLAGSHIP_CONFIG_FILE=""
ENV FLAGSHIP_ENV_ID=""
ENV FLAGSHIP_POLLING_INTERVAL=2000
ENV FLAGSHIP_BUCKETING_DIRECTORY=flagship

COPY --from=builder /dist/app .

# Command to run
CMD ["./app"]