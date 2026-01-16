package mutagens

import (
	"go/ast"
	"go/token"

	m "github.com/mouse-blink/gooze/internal/model"
)

func GenerateArithmeticMutationsV2(n ast.Node, fset *token.FileSet, content []byte, source m.SourceV2, mutationID *int) []m.MutationV2 {
	binExpr, ok := n.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	if !isArithmeticOpV2(binExpr.Op) {
		return nil
	}

	start, ok := offsetForPosV2(fset, binExpr.OpPos)
	if !ok {
		return nil
	}

	original := binExpr.Op.String()
	end := start + len(original)

	var mutations []m.MutationV2
	for _, mutatedOp := range getArithmeticAlternativesV2(binExpr.Op) {
		*mutationID++
		mutatedCode := replaceRangeV2(content, start, end, mutatedOp.String())
		mutations = append(mutations, m.MutationV2{
			ID:          uint(*mutationID - 1),
			Source:      source,
			Type:        m.MutationArithmetic,
			MutatedCode: mutatedCode,
		})
	}

	return mutations
}

func isArithmeticOpV2(op token.Token) bool {
	return op == token.ADD || op == token.SUB || op == token.MUL || op == token.QUO || op == token.REM
}

func getArithmeticAlternativesV2(original token.Token) []token.Token {
	allOps := []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM}

	var alternatives []token.Token
	for _, op := range allOps {
		if op != original {
			alternatives = append(alternatives, op)
		}
	}

	return alternatives
}
