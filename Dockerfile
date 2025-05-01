FROM golang:1.24-alpine

# Set working directory
WORKDIR /workspace

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Create directories
RUN mkdir -p /inputs
RUN mkdir -p /outputs

# Copy go script
COPY main.go .

# Set entrypoint
ENTRYPOINT ["go", "run", "main.go"]
