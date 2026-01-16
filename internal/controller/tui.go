package controller

import (
	"fmt"
	"io"

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

// Start initializes the UI.
func (t *TUI) Start() error {
	return nil
}

// Close finalizes the UI.
func (t *TUI) Close() {

}

// DisplayEstimation prints the estimation results or error.
func (t *TUI) DisplayEstimation(mutations []m.Mutation, err error) error {
	if err != nil {
		_, _ = fmt.Fprintf(t.output, "estimation error: %v\n", err)

		return err
	}

	paths := make(map[string]struct{})

	for _, mutation := range mutations {
		if mutation.Source.Origin == nil {
			continue
		}

		paths[string(mutation.Source.Origin.Path)] = struct{}{}
	}

	_, _ = fmt.Fprintf(t.output, "Estimated %d mutations across %d paths\n", len(mutations), len(paths))

	return nil
}
