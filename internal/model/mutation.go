// Package model defines the data structures for mutation testing.
package model

import "go/token"

// MutationType represents the category of mutation.
type MutationType string

const (
	// MutationArithmetic represents arithmetic operator mutations (+, -, *, /, %).
	MutationArithmetic MutationType = "arithmetic"
)

// Mutation represents a code mutation for testing.
type Mutation struct {
	ID         string
	Type       MutationType
	SourceFile Path
	OriginalOp token.Token
	MutatedOp  token.Token
	Line       int
	Column     int
	ScopeType  ScopeType
}
