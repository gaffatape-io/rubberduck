package testing

import (
	"os"
	"os/exec"
	"testing"
)

func FindBinary(t *testing.T) string {
	const bin = "rubberduck"
	
	// Check that the rubberduck binary is built, without it this test will fail.

	if path, err := exec.LookPath(bin); err == nil {
		return path
	}
	
	if _, err := os.Stat(bin); err != nil {
		t.Fatal("rubberduck not built; run > go build .")
	}

	return bin
}
