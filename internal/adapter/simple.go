package adapter

import (
	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/spf13/cobra"
)

// SimpleUI implements UI using cobra Command's Println.
type SimpleUI struct {
	cmd *cobra.Command
}

// NewSimpleUI creates a new SimpleUI.
func NewSimpleUI(cmd *cobra.Command) *SimpleUI {
	return &SimpleUI{cmd: cmd}
}

// Display shows source files using simple text output.
func (p *SimpleUI) Display(sources []m.Source) error {
	if len(sources) == 0 {
		p.cmd.Println("No source files found")
		return nil
	}

	for _, source := range sources {
		p.cmd.Println(source.Origin)
	}

	return nil
}

// ShowNotImplemented displays a "not implemented" message.
func (p *SimpleUI) ShowNotImplemented(count int) error {
	p.cmd.Printf("Found %d source files\n", count)
	p.cmd.Println("Mutation testing not yet implemented. Use --list to see source files.")

	return nil
}
