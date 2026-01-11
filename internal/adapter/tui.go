package adapter

import (
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	m "github.com/mouse-blink/gooze/internal/model"
)

// TUI implements UI using Bubble Tea for interactive display.
type TUI struct {
	output io.Writer
}

// NewTUI creates a new TUI.
func NewTUI(output io.Writer) *TUI {
	return &TUI{output: output}
}

// Display shows source files using Bubble Tea TUI.
func (p *TUI) Display(sources []m.Source) error {
	model := newTUIModel(sources, false)
	
	program := tea.NewProgram(model, tea.WithOutput(p.output))
	if _, err := program.Run(); err != nil {
		return err
	}
	
	return nil
}

// ShowNotImplemented displays a "not implemented" message using TUI.
func (p *TUI) ShowNotImplemented(count int) error {
	model := newNotImplementedModel(count)
	
	program := tea.NewProgram(model, tea.WithOutput(p.output))
	if _, err := program.Run(); err != nil {
		return err
	}
	
	return nil
}

// tuiModel represents the Bubble Tea model for displaying source files.
type tuiModel struct {
	sources        []m.Source
	notImplemented bool
	count          int
}

func newTUIModel(sources []m.Source, notImplemented bool) tuiModel {
	return tuiModel{
		sources:        sources,
		notImplemented: notImplemented,
	}
}

func newNotImplementedModel(count int) tuiModel {
	return tuiModel{
		notImplemented: true,
		count:          count,
	}
}

func (tm tuiModel) Init() tea.Cmd {
	return tea.Quit // Immediately quit after rendering
}

func (tm tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return tm, tea.Quit
}

func (tm tuiModel) View() string {
	var b strings.Builder
	
	// Header
	b.WriteString("╔════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║                    Gooze - Mutation Testing                    ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")
	
	if tm.notImplemented {
		// Show "not implemented" message
		b.WriteString(fmt.Sprintf("  📊 Found %d source file(s)\n\n", tm.count))
		b.WriteString("  ⚠️  Mutation testing not yet implemented.\n")
		b.WriteString("  💡 Use --list flag to see source files.\n\n")
		return b.String()
	}
	
	if len(tm.sources) == 0 {
		b.WriteString("  📭 No source files found\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("  📁 Found %d source file(s):\n\n", len(tm.sources)))
	
	for i, source := range tm.sources {
		b.WriteString(fmt.Sprintf("  %2d. %s\n", i+1, source.Origin))
	}
	
	b.WriteString("\n")
	return b.String()
}
