//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dev mg.Namespace
type Build mg.Namespace
type Test mg.Namespace
type Lint mg.Namespace
type Generate mg.Namespace
type Docs mg.Namespace

// Default target to run when none is specified
func Default() {
	fmt.Println("Available targets:")
	fmt.Println("  mage run <command>     - Run a CLI command (e.g., mage run pack)")
	fmt.Println("  mage dev:run <command> - Run a CLI command (same as mage run)")
	fmt.Println("  mage dev:server        - Run the development server")
	fmt.Println("  mage build             - Build the unified binary")
	fmt.Println("  mage test              - Run tests")
	fmt.Println("  mage lint              - Run linters")
	fmt.Println("  mage generate          - Generate code from OpenAPI spec")
	fmt.Println("  mage docs:gen          - Generate Swagger documentation")
}

// Run a CLI command (alias for dev:run)
func Run(what string) error {
	return Dev{}.Run(what)
}

// Run a CLI command
func (Dev) Run(what string) error {
	output, err := sh.Output("go", "run", "cmd/titanium/main.go", what)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}

// Run the development server
func (Dev) Server() error {
	return sh.Run("go", "run", "cmd/titanium/main.go", "--mode=server")
}

// Build the unified binary
func (Build) All() error {
	if err := os.MkdirAll("bin", 0755); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "bin/ti", "cmd/titanium/main.go")
}

// Run tests
func (Test) All() error {
	return sh.Run("go", "test", "./...")
}

// Run linters
func (Lint) All() error {
	if err := sh.Run("go", "vet", "./..."); err != nil {
		return err
	}
	return sh.Run("golangci-lint", "run")
}

// Generate code from OpenAPI spec
func (Generate) All() error {
	cmd := exec.Command("oapi-codegen", "-generate", "types,server", "-package", "api", "api/spec/project.yaml")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	return os.WriteFile("api/generated.go", output, 0644)
}

// Generate Swagger documentation
func (Docs) Gen() error {
	return sh.Run("swag", "init", "-g", "cmd/titanium/main.go", "-o", "docs")
}

// Clean build artifacts
func Clean() error {
	if err := os.RemoveAll("bin"); err != nil {
		return err
	}
	return sh.Run("go", "clean")
}
