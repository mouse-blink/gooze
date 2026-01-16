// Package cmd provides the root command and CLI setup for gooze.
package cmd

import (
	"os"

	"github.com/mouse-blink/gooze/internal/adapter"
	"github.com/mouse-blink/gooze/internal/domain"
	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/spf13/cobra"
)

var soirceFSAdapter adapter.SourceFSAdapter
var reportStore adapter.ReportStore
var fsAdapter adapter.SourceFSAdapter
var testAdapter adapter.TestRunnerAdapter
var orchestrator domain.Orchestrator
var workflow domain.WorkflowV2

func init() {
	soirceFSAdapter = adapter.NewLocalSourceFSAdapter()
	reportStore = adapter.NewReportStore()
	fsAdapter = adapter.NewLocalSourceFSAdapter()
	testAdapter = adapter.NewLocalTestRunnerAdapter()
	orchestrator = domain.NewOrchestrator(fsAdapter, testAdapter)
	workflow = domain.NewWorkflowV2(
		soirceFSAdapter,
		reportStore,
		orchestrator,
		domain.NewMutagen(),
	)
}

var listFlag bool
var parallelFlag int

// rootCmd represents the base command when called without any subcommands.
var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gooze [paths...]",
		Short: "Go mutation testing tool",
		Long: `Gooze is a mutation testing tool for Go that helps you assess the quality
of your test suite by introducing small changes (mutations) to your code
and verifying that your tests catch them.

Supports Go-style path patterns:
  - ./...          recursively scan current directory
  - ./pkg/...      recursively scan pkg directory
  - ./cmd ./pkg    scan multiple directories`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var paths []m.Path
			for _, arg := range args {
				paths = append(paths, m.Path(arg))
			}

			estimateArgs := domain.EstimateArgs{
							Paths:    paths,
							UseCache: listFlag,
						}
			if listFlag {
				return workflow.Estimate(estimateArgs)
			}
			return workflow.Test(domain.TestArgs{
				EstimateArgs: estimateArgs,
				Reports: ".gooze-reports",
				Threads: uint(parallelFlag),
			})
		},
	}
	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "list all source files and count of mutations applicable")
	cmd.Flags().IntVarP(&parallelFlag, "parallel", "p", 1, "number of parallel workers for mutation testing")
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
