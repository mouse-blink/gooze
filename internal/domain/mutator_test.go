package domain

import (
	"go/token"
	"path/filepath"
	"testing"

	m "github.com/mouse-blink/gooze/internal/model"
)

func TestGenerateMutations(t *testing.T) {
	t.Run("generates arithmetic mutations for addition operator in examples/basic", func(t *testing.T) {
		// Use examples/basic/main.go which has 3+5 expression
		basicPath := filepath.Join("..", "..", "examples", "basic", "main.go")
		source := loadSourceFromFile(t, basicPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		// Should generate 4 mutations: + → -, *, /, %
		if len(mutations) != 4 {
			t.Fatalf("expected 4 mutations for +, got %d", len(mutations))
		}

		// Verify all mutations are arithmetic type
		for _, mut := range mutations {
			if mut.Type != m.MutationArithmetic {
				t.Errorf("expected arithmetic mutation, got %v", mut.Type)
			}
			if mut.OriginalOp != token.ADD {
				t.Errorf("expected original op '+', got %s", mut.OriginalOp)
			}
			if mut.ScopeType != m.ScopeFunction {
				t.Errorf("expected function scope, got %v", mut.ScopeType)
			}
		}

		// Check we have all expected operators
		expectedOps := map[token.Token]bool{token.SUB: false, token.MUL: false, token.QUO: false, token.REM: false}
		for _, mut := range mutations {
			if _, ok := expectedOps[mut.MutatedOp]; ok {
				expectedOps[mut.MutatedOp] = true
			}
		}
		for op, found := range expectedOps {
			if !found {
				t.Errorf("missing mutation for operator: %s", op)
			}
		}
	})

	t.Run("generates mutations for arithmetic in examples/scopes", func(t *testing.T) {
		// Use examples/scopes/main.go which has + and - operators
		scopesPath := filepath.Join("..", "..", "examples", "scopes", "main.go")
		source := loadSourceFromFile(t, scopesPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		// scopes/main.go has: a + b and a - b in Calculate function
		// That's 2 operators × 4 mutations each = 8 mutations
		if len(mutations) < 8 {
			t.Fatalf("expected at least 8 mutations, got %d", len(mutations))
		}

		// Count mutations by original operator
		opCounts := make(map[token.Token]int)
		for _, mut := range mutations {
			opCounts[mut.OriginalOp]++
		}

		// Should have mutations for + and -
		if opCounts[token.ADD] == 0 {
			t.Error("expected mutations for + operator")
		}
		if opCounts[token.SUB] == 0 {
			t.Error("expected mutations for - operator")
		}
	})

	t.Run("assigns correct scope types to mutations", func(t *testing.T) {
		// examples/scopes/main.go has global consts, init func, and regular functions
		scopesPath := filepath.Join("..", "..", "examples", "scopes", "main.go")
		source := loadSourceFromFile(t, scopesPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		// Count by scope type
		scopeCounts := make(map[m.ScopeType]int)
		for _, mut := range mutations {
			scopeCounts[mut.ScopeType]++
		}

		// Should have mutations in function scope (Calculate function)
		if scopeCounts[m.ScopeFunction] == 0 {
			t.Error("expected mutations in function scope")
		}
	})

	t.Run("no mutations when no arithmetic operators present", func(t *testing.T) {
		// examples/empty/main.go has comparison operators but no arithmetic
		emptyPath := filepath.Join("..", "..", "examples", "empty", "main.go")
		source := loadSourceFromFile(t, emptyPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		if len(mutations) != 0 {
			t.Fatalf("expected 0 mutations for file without arithmetic, got %d", len(mutations))
		}
	})

	t.Run("tracks line and column positions correctly", func(t *testing.T) {
		basicPath := filepath.Join("..", "..", "examples", "basic", "main.go")
		source := loadSourceFromFile(t, basicPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		for _, mut := range mutations {
			if mut.Line == 0 {
				t.Error("mutation Line should not be 0")
			}
			if mut.Column == 0 {
				t.Error("mutation Column should not be 0")
			}
		}
	})

	t.Run("generates unique IDs for each mutation", func(t *testing.T) {
		basicPath := filepath.Join("..", "..", "examples", "basic", "main.go")
		source := loadSourceFromFile(t, basicPath)

		wf := NewWorkflow()
		mutations, err := wf.GenerateMutations(source)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		ids := make(map[string]bool)
		for _, mut := range mutations {
			if mut.ID == "" {
				t.Error("mutation ID should not be empty")
			}
			if ids[mut.ID] {
				t.Errorf("duplicate mutation ID: %s", mut.ID)
			}
			ids[mut.ID] = true
		}
	})
}

// Helper function to load a source from an actual file.
func loadSourceFromFile(t *testing.T, path string) m.Source {
	t.Helper()

	// Use the workflow's existing GetSources method
	wf := NewWorkflow()
	sources, err := wf.GetSources(m.Path(path))
	if err != nil {
		t.Fatalf("failed to load source from %s: %v", path, err)
	}

	if len(sources) == 0 {
		t.Fatalf("no sources found in %s", path)
	}

	return sources[0]
}
