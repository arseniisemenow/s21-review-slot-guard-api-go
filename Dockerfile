FROM golang:1.23-bookworm

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y make git

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Keep container running
CMD ["sleep", "infinity"]
