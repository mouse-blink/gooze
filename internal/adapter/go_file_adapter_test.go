package adapter

import (
	"path/filepath"
	"testing"

	"go/token"
)

func TestLocalGoFileAdapter_Parse(t *testing.T) {
	adapter := NewLocalGoFileAdapter()
	fset := token.NewFileSet()

	exampleFile := filepath.Join(examplePath(t, "basic"), "main.go")
	content := readFileBytes(t, exampleFile)
	file, err := adapter.Parse(fset, exampleFile, content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if file.Name.Name != "main" {
		t.Fatalf("Parse() package = %s, want main", file.Name.Name)
	}
}

func TestLocalGoFileAdapter_Parse_InvalidSource(t *testing.T) {
	adapter := NewLocalGoFileAdapter()
	fset := token.NewFileSet()

	if _, err := adapter.Parse(fset, "broken.go", []byte("package foo\n func")); err == nil {
		t.Fatalf("Parse() expected error for invalid source")
	}
}
