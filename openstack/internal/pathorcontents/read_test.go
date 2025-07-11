package pathorcontents

import (
	"io"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestUnitRead_Path(t *testing.T) {
	isPath := true

	f, cleanup := testTempFile(t)
	defer cleanup()

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	f.Close()

	contents, wasPath, err := Read(f.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}

	if contents != "foobar" {
		t.Fatalf("expected contents %s, got %s", "foobar", contents)
	}
}

func TestUnitRead_TildePath(t *testing.T) {
	isPath := true

	home, err := homedir.Dir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	f, cleanup := testTempFile(t, home)
	defer cleanup()

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	f.Close()

	r := strings.NewReplacer(home, "~")
	homePath := r.Replace(f.Name())

	contents, wasPath, err := Read(homePath)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}

	if contents != "foobar" {
		t.Fatalf("expected contents %s, got %s", "foobar", contents)
	}
}

func TestUnitRead_PathNoPermission(t *testing.T) {
	// This skip condition is intended to get this test out of the way of users
	// who are building and testing Terraform from within a Linux-based Docker
	// container, where it is common for processes to be running as effectively
	// root within the container.
	if u, err := user.Current(); err == nil && u.Uid == "0" {
		t.Skip("This test is invalid when running as root, since root can read every file")
	}

	isPath := true

	f, cleanup := testTempFile(t)
	defer cleanup()

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	f.Close()

	if err := os.Chmod(f.Name(), 0); err != nil {
		t.Fatalf("err: %s", err)
	}

	contents, wasPath, err := Read(f.Name())

	if err == nil {
		t.Fatal("Expected error, got none!")
	}

	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}

	if contents != "" {
		t.Fatalf("expected contents %s, got %s", "", contents)
	}
}

func TestUnitRead_Contents(t *testing.T) {
	isPath := false
	input := "hello"

	contents, wasPath, err := Read(input)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}

	if contents != input {
		t.Fatalf("expected contents %s, got %s", input, contents)
	}
}

func TestUnitRead_TildeContents(t *testing.T) {
	isPath := false
	input := "~/hello/notafile"

	contents, wasPath, err := Read(input)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}

	if contents != input {
		t.Fatalf("expected contents %s, got %s", input, contents)
	}
}

// Returns an open tempfile based at baseDir and a function to clean it up.
func testTempFile(t *testing.T, baseDir ...string) (*os.File, func()) {
	base := ""
	if len(baseDir) == 1 {
		base = baseDir[0]
	}

	f, err := os.CreateTemp(base, "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return f, func() {
		os.Remove(f.Name())
	}
}
