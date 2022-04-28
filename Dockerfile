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

ENV FS_ENV_ID=""
ENV FS_POLLING_INTERVAL=2000
ENV FS_PORT=8080
ENV FS_ADDRESS="0.0.0.0"

COPY --from=builder /dist/app .

EXPOSE ${FS_PORT}

# Command to run
CMD ["sh","-c","./app --envId=${FS_ENV_ID} --pollingInterval=${FS_POLLING_INTERVAL} --port=${FS_PORT} --address=${FS_ADDRESS}"]