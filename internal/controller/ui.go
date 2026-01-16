// Package controller provides output adapters for displaying mutation testing results.
package controller

import (
	m "github.com/mouse-blink/gooze/internal/model"
)

// MutationEstimation holds estimation counts for different mutation types.
type MutationEstimation struct {
	Arithmetic int
	Boolean    int
}

// UI defines the interface for displaying source file lists.
// Implementations can use different output methods (simple text, TUI, etc).
type UI interface {
	Start() error
	Close()
	DisplayEstimation(mutations []m.Mutation, err error) error
	DisplayConcurencyInfo(threads int, count int)
	DusplayUpcomingTestsInfo(i int)
	DisplayStartingTestInfo(currentMutation m.Mutation)
	DisplayCompletedTestInfo(currentMutation m.Mutation, mutationResult m.Result)
}
