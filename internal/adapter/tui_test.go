package adapter

import (
	"bytes"
	"strings"
	"testing"

	m "github.com/mouse-blink/gooze/internal/model"
)

func TestTUI_Display(t *testing.T) {
	tests := []struct {
		name          string
		sources       []m.Source
		wantContains  []string
		wantNotEmpty  bool
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
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Create TUI
			ui := NewTUI(&buf)
			err := ui.Display(tt.sources)

			if err != nil {
				t.Errorf("Display() error = %v", err)
				return
			}

			got := buf.String()

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
