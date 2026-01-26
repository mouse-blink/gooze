package cmd

import (
	"github.com/mouse-blink/gooze/internal/domain"
	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command.
var viewCmd = newViewCmd()

func newViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [reports-dir]",
		Short: "View previously generated mutation reports",
		Long:  "View previously generated mutation reports from a reports directory.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			dir := reportsOutputDirFlag
			if len(args) == 1 {
				dir = args[0]
			}

			return workflow.View(domain.ViewArgs{Reports: m.Path(dir)})
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
