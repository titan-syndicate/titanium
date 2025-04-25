# Titanium

A Go-based platform API for managing projects and namespaces.

## Prerequisites

- Go 1.21 or later
- Mage (build tool)

## Installation

1. Install Mage:
   ```bash
   brew install mage
   ```

2. Install project dependencies:
   ```bash
   go mod download
   ```

## Development

This project uses Mage as its build tool. Available commands:

```bash
# Show available commands
mage

# Run the development server
mage dev:run

# Build the server binary
mage build

# Run tests
mage test

# Run linters
mage lint

# Generate code from OpenAPI spec
mage generate

# Clean build artifacts
mage clean
```

## API Documentation

The API documentation is served using an embedded Swagger UI. Once the server is running, you can access:
- Swagger UI: http://localhost:8080/swagger/
- OpenAPI spec: http://localhost:8080/api/spec/project.yaml