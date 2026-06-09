package main

import (
	"os"
	"testing"
)

func TestMainPackage_StartsWithoutServer(t *testing.T) {
	origRunServer := runServer
	runServer = func(params ...string) {}
	defer func() { runServer = origRunServer }()

	main()
}

func TestMainPackage_LogsWhenDotenvMissing(t *testing.T) {
	origRunServer := runServer
	runServer = func(params ...string) {}
	defer func() { runServer = origRunServer }()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	main()
}
