package adapter

import (
	"bytes"
	"testing"

	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/spf13/cobra"
)

func TestSimpleUI_Display(t *testing.T) {
	tests := []struct {
		name    string
		sources []m.Source
		want    string
	}{
		{
			name:    "empty list",
			sources: []m.Source{},
			want:    "No source files found\n",
		},
		{
			name: "single file",
			sources: []m.Source{
				{Origin: m.Path("main.go")},
			},
			want: "main.go\n",
		},
		{
			name: "multiple files",
			sources: []m.Source{
				{Origin: m.Path("main.go")},
				{Origin: m.Path("helper.go")},
				{Origin: m.Path("types.go")},
			},
			want: "main.go\nhelper.go\ntypes.go\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a cobra command with a buffer to capture output
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&buf)

			// Create UI and display list
			ui := NewSimpleUI(cmd)
			err := ui.Display(tt.sources)

			if err != nil {
				t.Errorf("Display() error = %v", err)
				return
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("Display() output = %q, want %q", got, tt.want)
			}
		})
	}
}
