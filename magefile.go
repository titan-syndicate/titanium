//go:build mage
// +build mage

package main

import (
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

// Run namespace for specific run commands
type Run mg.Namespace

// Build the node test app
func (Run) NodeTest() error {
	cmd := exec.Command("go", "run", "cmd/titanium/main.go", "build", "./test/node-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build the go test app
func (Run) GoTest() error {
	cmd := exec.Command("go", "run", "cmd/titanium/main.go", "build", "go", "v0.0.0-test", "./test/go-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	if err := os.RemoveAll("dist"); err != nil {
		return err
	}
	if err := os.RemoveAll(".tmp"); err != nil {
		return err
	}
	return sh.Run("go", "clean")
}
