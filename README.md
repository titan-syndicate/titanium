# Titanium

A Go-based platform API for managing projects and namespaces.

## Prerequisites

- Go 1.21 or later
- Mage (build tool)
- Docker (for building and running containers)

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

# Build the unified binary
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

## Running the Application

### Using the Binary

After building the application with `mage build`, you can run commands directly using the binary:

```bash
# Build a Node.js application
./bin/ti build ./test/node-app

# Run the server
./bin/ti --mode=server
```

### Using Mage

For convenience, we provide mage commands for common operations:

```bash
# Build the Node.js test application
mage run:nodeTest

# Run the development server
mage dev:server
```

## Testing

### Node.js Test Application

The project includes a simple Node.js test application in `test/node-app/`. You can build it using either:

```bash
# Using the binary
./bin/ti build ./test/node-app

# Using mage
mage run:nodeTest
```

The build process will:
1. Use Cloud Native Buildpacks to detect the application type
2. Install Node.js and dependencies
3. Create a container image named `test-app`

### Running the Built Container

After building, you can run the container:

```bash
docker run -p 8080:8080 test-app
```

## API Documentation

The API documentation is served using an embedded Swagger UI. Once the server is running, you can access:
- Swagger UI: http://localhost:8080/swagger/
- OpenAPI spec: http://localhost:8080/api/spec/project.yaml