# Project Orca Docker Image

This document provides information about the official Docker image for Project Orca.

## Quick Start

Pull the latest image:
```bash
docker pull rianfowler/project-orca:latest
```

Run the container:
```bash
docker run --rm rianfowler/project-orca:latest ti --help
```

## Available Tags

- `latest`: Latest stable release
- `vX.Y.Z`: Specific version (e.g., `v1.0.0`)

## Multi-Architecture Support

The image is available for the following architectures:
- linux/amd64
- linux/arm64

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TITANIUM_CONFIG` | Path to configuration file | `/etc/titanium/config.yaml` |

## Volume Mounts

Common volume mounts:
```bash
# Mount a local config file
docker run --rm -v $(pwd)/config.yaml:/etc/titanium/config.yaml rianfowler/project-orca:latest

# Mount a data directory
docker run --rm -v $(pwd)/data:/data rianfowler/project-orca:latest
```

## Building Locally

To build the image locally:
```bash
docker build -t rianfowler/project-orca:local .
```

## Contributing

To contribute to the Docker image, please refer to the [Dockerfile](Dockerfile) in the repository root.