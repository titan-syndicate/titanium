FROM alpine:latest

# Install CA certificates (if your CLI needs HTTPS support)
RUN apk add --no-cache ca-certificates

# Copy the CLI binary into the container
COPY demp /usr/local/bin/demp

# Make sure the binary is executable
RUN chmod +x /usr/local/bin/demp

# Set the entrypoint so that any container arguments are passed to the CLI
ENTRYPOINT ["/usr/local/bin/ti"]