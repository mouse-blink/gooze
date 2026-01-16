package controller

import (
	"bytes"
	"fmt"
	"sort"

	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/olekukonko/tablewriter"
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

// Start initializes the UI.
func (s *SimpleUI) Start() error {
	return nil
}

// Close finalizes the UI.
func (s *SimpleUI) Close() {

}

// DisplayEstimation prints the estimation results or error.
func (s *SimpleUI) DisplayEstimation(mutations []m.Mutation, err error) error {
	if err != nil {
		s.printf("estimation error: %v\n", err)
		return err
	}

	info := make(map[string]int)

	for _, mutation := range mutations {
		if mutation.Source.Origin == nil {
			continue
		}

		info[string(mutation.Source.Origin.Path)]++
	}

	pathsList := make([]string, 0, len(info))
	for path := range info {
		pathsList = append(pathsList, path)
	}

	sort.Strings(pathsList)

	var tableBuffer bytes.Buffer

	table := tablewriter.NewWriter(&tableBuffer)
	table.SetHeader([]string{"Path", "Mutations"})
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

	pathsCount := 0

	for _, pathStr := range pathsList {
		count := info[pathStr]
		table.Append([]string{pathStr, fmt.Sprintf("%d", count)})

		pathsCount++
	}

	table.SetFooter([]string{
		fmt.Sprintf("Total Files %d", pathsCount),
		fmt.Sprintf("%d", len(mutations)),
	})

	table.Render()
	s.printf("\n%s", tableBuffer.String())

	return nil
}

func (s *SimpleUI) printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(s.cmd.OutOrStdout(), format, args...)
}
