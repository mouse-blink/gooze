package mutagens

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestOffsetForPos(t *testing.T) {
	booleanPath := filepath.Join("..", "..", "..", "examples", "boolean", "main.go")
	content, err := os.ReadFile(booleanPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", booleanPath, err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, booleanPath, content, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse source %s: %v", booleanPath, err)
	}

	var literalPos token.Pos
	ast.Inspect(file, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == "true" || ident.Name == "false" {
			literalPos = ident.Pos()
			return false
		}
		return true
	})

	if literalPos == token.NoPos {
		t.Fatal("expected to find boolean literal")
	}

	offset, ok := offsetForPos(fset, literalPos)
	if !ok {
		t.Fatal("expected offsetForPosV2 to return ok")
	}

	if !(bytes.HasPrefix(content[offset:], []byte("true")) || bytes.HasPrefix(content[offset:], []byte("false"))) {
		start := offset
		end := offset + 10
		if start < 0 {
			start = 0
		}
		if end > len(content) {
			end = len(content)
		}
		snippet := string(content[start:end])
		t.Fatalf("unexpected offset %d; snippet=%q", offset, snippet)
	}
}

func TestReplaceRange(t *testing.T) {
	original := []byte("abcde")
	mutated := replaceRange(original, 1, 4, "XYZ")

	if string(mutated) != "aXYZe" {
		t.Fatalf("replaceRangeV2 result = %q, expected %q", string(mutated), "aXYZe")
	}

	unchanged := replaceRange(original, 10, 12, "nope")
	if !bytes.Equal(unchanged, original) {
		t.Fatalf("expected out-of-range replacement to return original content")
	}
}

func TestDiffCode(t *testing.T) {
	original := []byte("package main\n\nfunc foo() {\n\treturn 3 + 5\n}")
	mutated := []byte("package main\n\nfunc foo() {\n\treturn 3 - 5\n}")

	diff := diffCode(original, mutated)
	diffStr := string(diff)

	if !bytes.Contains(diff, []byte("--- original")) {
		t.Errorf("expected diff to contain '--- original', got: %s", diffStr)
	}
	if !bytes.Contains(diff, []byte("+++ mutated")) {
		t.Errorf("expected diff to contain '+++ mutated', got: %s", diffStr)
	}
	if !bytes.Contains(diff, []byte("-\treturn 3 + 5")) {
		t.Errorf("expected diff to contain removed line, got: %s", diffStr)
	}
	if !bytes.Contains(diff, []byte("+\treturn 3 - 5")) {
		t.Errorf("expected diff to contain added line, got: %s", diffStr)
	}
}

func TestDiffCode_EmptyInputs(t *testing.T) {
	// Test with both empty
	diff := diffCode(nil, nil)
	if diff != nil {
		t.Errorf("expected nil diff for empty inputs, got: %v", diff)
	}

	diff = diffCode([]byte{}, []byte{})
	if diff != nil {
		t.Errorf("expected nil diff for empty byte slices, got: %v", diff)
	}

	// Test with one empty
	original := []byte("some content")
	diff = diffCode(original, nil)
	if diff == nil {
		t.Errorf("expected diff when one input is empty")
	}

	diff = diffCode(nil, original)
	if diff == nil {
		t.Errorf("expected diff when one input is empty")
	}
}

func TestEnsureTrailingNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "empty content",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "content with newline",
			input:    []byte("hello world\n"),
			expected: []byte("hello world\n"),
		},
		{
			name:     "content without newline",
			input:    []byte("hello world"),
			expected: []byte("hello world\n"),
		},
		{
			name:     "multiline with newline",
			input:    []byte("line1\nline2\n"),
			expected: []byte("line1\nline2\n"),
		},
		{
			name:     "multiline without newline",
			input:    []byte("line1\nline2"),
			expected: []byte("line1\nline2\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureTrailingNewline(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ensureTrailingNewline() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
