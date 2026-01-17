package mutagens

import (
	"go/ast"
	"go/token"

	m "github.com/mouse-blink/gooze/internal/model"
)

// GenerateLogicalMutations generates logical operator mutations for the given AST node.
func GenerateLogicalMutations(n ast.Node, fset *token.FileSet, content []byte, source m.Source, mutationID *int) []m.Mutation {
	binExpr, ok := n.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	if !isLogicalOp(binExpr.Op) {
		return nil
	}

	start, ok := offsetForPos(fset, binExpr.OpPos)
	if !ok {
		return nil
	}

	original := binExpr.Op.String()
	end := start + len(original)

	var mutations []m.Mutation

	for _, mutatedOp := range getLogicalAlternatives(binExpr.Op) {
		*mutationID++
		mutatedCode := replaceRange(content, start, end, mutatedOp.String())
		mutations = append(mutations, m.Mutation{
			ID:          *mutationID - 1,
			Source:      source,
			Type:        m.MutationLogical,
			MutatedCode: mutatedCode,
		})
	}

	return mutations
}

func isLogicalOp(op token.Token) bool {
	return op == token.LAND || op == token.LOR
}

func getLogicalAlternatives(original token.Token) []token.Token {
	allOps := []token.Token{token.LAND, token.LOR}

	var alternatives []token.Token

	for _, op := range allOps {
		if op != original {
			alternatives = append(alternatives, op)
		}
	}

	return alternatives
}
