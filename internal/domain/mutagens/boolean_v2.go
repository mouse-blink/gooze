package mutagens

import (
	"go/ast"
	"go/token"

	m "github.com/mouse-blink/gooze/internal/model"
)

func GenerateBooleanMutationsV2(n ast.Node, fset *token.FileSet, content []byte, source m.SourceV2, mutationID *int) []m.MutationV2 {
	ident, ok := n.(*ast.Ident)
	if !ok {
		return nil
	}

	if !isBooleanLiteralV2(ident.Name) {
		return nil
	}

	start, ok := offsetForPosV2(fset, ident.Pos())
	if !ok {
		return nil
	}

	end := start + len(ident.Name)
	mutated := flipBooleanV2(ident.Name)

	*mutationID++
	mutatedCode := replaceRangeV2(content, start, end, mutated)
	return []m.MutationV2{{
		ID:          uint(*mutationID - 1),
		Source:      source,
		Type:        m.MutationBoolean,
		MutatedCode: mutatedCode,
	}}
}

func isBooleanLiteralV2(name string) bool {
	return name == "true" || name == "false"
}

func flipBooleanV2(original string) string {
	if original == "true" {
		return "false"
	}

	return "true"
}
