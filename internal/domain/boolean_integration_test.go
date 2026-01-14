package domain

import (
	"path/filepath"
	"testing"

	"github.com/mouse-blink/gooze/internal/adapter"

	m "github.com/mouse-blink/gooze/internal/model"
)

func TestBooleanMutationIntegration(t *testing.T) {
	t.Run("end-to-end boolean mutation kills test", func(t *testing.T) {
		booleanPath := filepath.Join("..", "..", "examples", "boolean", "main.go")
		source := loadSourceFromFile(t, booleanPath)

		wf := NewWorkflow(adapter.NewLocalSourceFSAdapter(), adapter.NewLocalGoFileAdapter(), adapter.NewLocalTestRunnerAdapter())

		// Generate boolean mutations
		mutations, err := wf.GenerateMutations(source, m.MutationBoolean)
		if err != nil {
			t.Fatalf("GenerateMutations failed: %v", err)
		}

		if len(mutations) == 0 {
			t.Fatal("expected at least one boolean mutation")
		}

		// Find a mutation that flips true to false in checkStatus function
		var testMutation *m.Mutation
		for i := range mutations {
			mut := &mutations[i]
			if mut.OriginalText == "true" && mut.MutatedText == "false" {
				testMutation = mut
				break
			}
		}

		if testMutation == nil {
			t.Fatal("could not find a true->false mutation")
		}

		// Test the mutation using RunMutationTests
		results, err := wf.RunMutationTests([]m.Source{source}, 1)
		if err != nil {
			t.Fatalf("RunMutationTests failed: %v", err)
		}

		fileResult, ok := results[source.Origin]
		if !ok {
			t.Fatal("no results for source")
		}

		// Find the report for our test mutation
		var report *m.Report
		for i := range fileResult.Reports {
			if fileResult.Reports[i].MutationID == testMutation.ID {
				report = &fileResult.Reports[i]
				break
			}
		}

		if report == nil {
			t.Fatalf("no report found for mutation %s", testMutation.ID)
		}

		// The test should kill the mutation (detect the change)
		if !report.Killed {
			t.Errorf("expected mutation to be killed by tests, but it survived")
			t.Logf("Mutation: %s at line %d, col %d", testMutation.ID, testMutation.Line, testMutation.Column)
			t.Logf("Report output: %s", report.Output)
		}
	})
}
