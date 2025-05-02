FROM alpine:latest

# Install CA certificates (if your CLI needs HTTPS support)
RUN apk add --no-cache

# Copy the CLI binary into the container
COPY ti /usr/local/bin/ti

# Make sure the binary is executable
RUN chmod +x /usr/local/bin/ti

# Set the entrypoint so that any container arguments are passed to the CLI
ENTRYPOINT ["/usr/local/bin/ti"]