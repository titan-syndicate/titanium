# Titanium Plugin Host

This package implements the plugin host system for Titanium. It is responsible for loading, managing, and executing plugins.

## Directory Structure

```
internal/pluginhost/
├── README.md           # This file
├── manager.go         # Core plugin management logic
├── proto/            # Protocol definitions and generated code
│   ├── plugin.proto  # Protocol definitions
│   └── gen/         # Generated code
└── example/         # Example plugin implementation
    └── plugin.go    # Example plugin
```

## Components

### Manager

The `manager.go` file contains the core plugin management logic:
- Plugin discovery
- Plugin loading
- Plugin execution
- Plugin lifecycle management

### Protocol

The `proto` directory contains:
- Protocol definitions for plugin communication
- Generated code for the host implementation

### Example

The `example` directory contains an example plugin implementation that demonstrates:
- How to implement the plugin interface
- How to use the gRPC server
- How to handle plugin lifecycle

## Usage

This package is internal to Titanium and should not be imported by plugins. Plugins should use the public interface defined in `pkg/plugin`.

## Development

When modifying this package:
1. Update the protocol definitions if needed
2. Regenerate the protocol code
3. Update the example plugin if the interface changes
4. Add tests for new functionality
5. Update this documentation