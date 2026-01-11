package adapter

import (
	"fmt"
	"strings"
	"testing"

	m "github.com/mouse-blink/gooze/internal/model"
)

func TestTUI_Display(t *testing.T) {
	tests := []struct {
		name         string
		sources      []m.Source
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name:    "empty list",
			sources: []m.Source{},
			wantContains: []string{
				"No source files found",
			},
			wantNotEmpty: true,
		},
		{
			name: "single file",
			sources: []m.Source{
				{Origin: m.Path("main.go")},
			},
			wantContains: []string{
				"main.go",
			},
			wantNotEmpty: true,
		},
		{
			name: "multiple files",
			sources: []m.Source{
				{Origin: m.Path("main.go")},
				{Origin: m.Path("helper.go")},
				{Origin: m.Path("types.go")},
			},
			wantContains: []string{
				"main.go",
				"helper.go",
				"types.go",
			},
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the View directly instead of full TUI
			model := newTUIModelWithHeight(tt.sources, false, 20)
			got := model.View()

			if tt.wantNotEmpty && got == "" {
				t.Errorf("Display() output is empty, want non-empty")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Display() output does not contain %q\nGot: %q", want, got)
				}
			}
		})
	}
}

func TestTUIModel_Pagination(t *testing.T) {
	t.Run("title always visible in view", func(t *testing.T) {
		sources := make([]m.Source, 50)
		for i := 0; i < 50; i++ {
			sources[i] = m.Source{Origin: m.Path(fmt.Sprintf("file%d.go", i))}
		}

		model := newTUIModelWithHeight(sources, false, 20)
		view := model.View()

		if !strings.Contains(view, "Gooze - Mutation Testing") {
			t.Errorf("Title not found in view")
		}
	})

	t.Run("shows pagination controls when list is large", func(t *testing.T) {
		sources := make([]m.Source, 50)
		for i := 0; i < 50; i++ {
			sources[i] = m.Source{Origin: m.Path(fmt.Sprintf("file%d.go", i))}
		}

		model := newTUIModelWithHeight(sources, false, 20)
		view := model.View()

		// Should show navigation hints when paginated
		hasNavigation := strings.Contains(view, "↑/↓") ||
			strings.Contains(view, "j/k") ||
			strings.Contains(view, "Page")
		if !hasNavigation {
			t.Errorf("Navigation hints not found in view with 50 items and height 20")
		}
	})

	t.Run("limits displayed items to fit screen", func(t *testing.T) {
		sources := make([]m.Source, 50)
		for i := 0; i < 50; i++ {
			sources[i] = m.Source{Origin: m.Path(fmt.Sprintf("file%d.go", i))}
		}

		model := newTUIModelWithHeight(sources, false, 15)
		view := model.View()

		// Count how many files are shown
		lines := strings.Split(view, "\n")
		fileCount := 0
		for _, line := range lines {
			if strings.Contains(line, ".go") {
				fileCount++
			}
		}

		// Should show fewer files than total due to pagination
		if fileCount >= 50 {
			t.Errorf("All 50 files shown, expected pagination. Got %d files", fileCount)
		}

		if fileCount == 0 {
			t.Errorf("No files shown, expected at least some files")
		}
	})

	t.Run("small list shows all items without pagination", func(t *testing.T) {
		sources := []m.Source{
			{Origin: m.Path("main.go")},
			{Origin: m.Path("helper.go")},
			{Origin: m.Path("types.go")},
		}

		model := newTUIModelWithHeight(sources, false, 20)
		view := model.View()

		// All 3 files should be visible
		if !strings.Contains(view, "main.go") {
			t.Errorf("main.go not found")
		}
		if !strings.Contains(view, "helper.go") {
			t.Errorf("helper.go not found")
		}
		if !strings.Contains(view, "types.go") {
			t.Errorf("types.go not found")
		}
	})
}

// Helper for tests to create model with specific height
func newTUIModelWithHeight(sources []m.Source, notImplemented bool, height int) tuiModel {
	model := tuiModel{
		sources:        sources,
		notImplemented: notImplemented,
		height:         height,
	}
	return model
}
