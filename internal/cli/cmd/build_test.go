package ti

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildCommand(t *testing.T) {
	// Get the absolute path to the test node app
	testAppPath, err := filepath.Abs("../../../../test/node-app")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Change to the test app directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(testAppPath); err != nil {
		t.Fatalf("Failed to change to test app directory: %v", err)
	}

	// Run the build command
	err = runBuild([]string{})
	if err != nil {
		t.Errorf("Build command failed: %v", err)
	}

	// TODO: Add assertions to verify the build output
	// For now, we'll just check that the command runs without error
}
