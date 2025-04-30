# Titanium Plugin System

This package provides the interface and tools needed to create plugins for Titanium.

## Overview

The plugin system uses gRPC for communication between the Titanium host and plugins. The interface is defined using Protocol Buffers, making it language-agnostic.

## Directory Structure

```
pkg/plugin/
├── README.md           # This file
└── proto/             # Protocol definitions
    ├── plugin.proto   # The plugin interface definition
    └── gen/          # Generated code for different languages
        └── go/       # Go-specific generated code
```

## Using the Plugin Interface

### Go Plugins

```go
import "github.com/titan-syndicate/titanium/pkg/plugin/proto/gen/go"

// Your plugin implementation
type MyPlugin struct {
    // ...
}

// Implement the required methods
func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Version() string {
    return "v1.0.0"
}

func (p *MyPlugin) Execute(args []string) (string, error) {
    // Your plugin logic here
    return "Hello from my plugin!", nil
}
```

### Other Languages

For other languages, you can either:
1. Use our pre-generated code (if available for your language)
2. Generate your own code from the proto definition

See the `proto` directory for the interface definition and instructions for generating code.

## Versioning

The plugin interface is versioned. Each version of Titanium will specify which version of the plugin interface it supports. Make sure your plugin implements the correct version.

## Contributing

When contributing to the plugin system:
1. Update the proto definition if adding new features
2. Generate code for all supported languages
3. Update this documentation
4. Add tests for new functionality