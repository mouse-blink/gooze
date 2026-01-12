package domain

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	m "github.com/mouse-blink/gooze/internal/model"
)

// GenerateMutations analyzes a source file and generates all arithmetic mutations.
func (w *workflow) GenerateMutations(source m.Source) ([]m.Mutation, error) {
	// Parse the source file
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, string(source.Origin), nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", source.Origin, err)
	}

	var mutations []m.Mutation

	mutationID := 0

	// Walk the AST and find arithmetic operators
	ast.Inspect(file, func(n ast.Node) bool {
		binExpr, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		// Check if it's an arithmetic operator
		if !isArithmeticOp(binExpr.Op) {
			return true
		}

		// Get position information
		pos := fset.Position(binExpr.OpPos)

		// Find which scope this mutation belongs to
		scopeType := findScopeType(source.Scopes, pos.Line)

		// Generate mutations for all alternative operators
		for _, mutatedOp := range getArithmeticAlternatives(binExpr.Op) {
			mutationID++
			mutations = append(mutations, m.Mutation{
				ID:         fmt.Sprintf("ARITH_%d", mutationID),
				Type:       m.MutationArithmetic,
				SourceFile: source.Origin,
				OriginalOp: binExpr.Op,
				MutatedOp:  mutatedOp,
				Line:       pos.Line,
				Column:     pos.Column,
				ScopeType:  scopeType,
			})
		}

		return true
	})

	return mutations, nil
}

// isArithmeticOp checks if a token is an arithmetic operator.
func isArithmeticOp(op token.Token) bool {
	return op == token.ADD || op == token.SUB || op == token.MUL || op == token.QUO || op == token.REM
}

// getArithmeticAlternatives returns all alternative operators for mutation.
func getArithmeticAlternatives(original token.Token) []token.Token {
	allOps := []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM}

	var alternatives []token.Token

	for _, op := range allOps {
		if op != original {
			alternatives = append(alternatives, op)
		}
	}

	return alternatives
}

// findScopeType determines which scope a line belongs to.
func findScopeType(scopes []m.CodeScope, line int) m.ScopeType {
	for _, scope := range scopes {
		if line >= scope.StartLine && line <= scope.EndLine {
			return scope.Type
		}
	}

	// Default to function scope if not found in any scope
	return m.ScopeFunction
}
