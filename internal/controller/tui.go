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

// DisplayConcurencyInfo shows concurrency settings.
func (t *TUI) DisplayConcurencyInfo(threads int, count int) {
	_, _ = fmt.Fprintf(t.output, "Running %d mutations with %d worker(s)\n", count, threads)
}

// DusplayUpcomingTestsInfo shows the number of upcoming mutations to be tested.
func (t *TUI) DusplayUpcomingTestsInfo(i int) {
	_, _ = fmt.Fprintf(t.output, "Upcoming mutations: %d\n", i)
}

// DisplayStartingTestInfo shows info about the mutation test starting.
func (t *TUI) DisplayStartingTestInfo(currentMutation m.Mutation) {
	path := ""
	if currentMutation.Source.Origin != nil {
		path = string(currentMutation.Source.Origin.Path)
	}

	_, _ = fmt.Fprintf(t.output, "Starting mutation %d (%s) %s\n", currentMutation.ID, currentMutation.Type, path)
}

// DisplayCompletedTestInfo shows info about the completed mutation test.
func (t *TUI) DisplayCompletedTestInfo(currentMutation m.Mutation, mutationResult m.Result) {
	status := "unknown"
	if results, ok := mutationResult[currentMutation.Type]; ok && len(results) > 0 {
		status = formatTestStatus(results[0].Status)
	}

	_, _ = fmt.Fprintf(t.output, "Completed mutation %d (%s) -> %s\n", currentMutation.ID, currentMutation.Type, status)
}
