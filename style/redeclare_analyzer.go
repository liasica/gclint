package style

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func newRedeclareAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "redeclare",
		Doc:  "report short variable declarations that reuse an existing variable in the same scope",
		Run:  runRedeclareAnalyzer,
	}
}

func runRedeclareAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			assignStatement, ok := node.(*ast.AssignStmt)
			if !ok || assignStatement.Tok != token.DEFINE {
				return true
			}

			if isTypeSwitchGuard(assignStatement) {
				return true
			}

			for _, leftHandSideExpression := range assignStatement.Lhs {
				identifier, ok := leftHandSideExpression.(*ast.Ident)
				if !ok || identifier.Name == "_" || identifier.Name == "err" {
					continue
				}

				if pass.TypesInfo.Defs[identifier] != nil {
					continue
				}

				pass.Reportf(
					identifier.Pos(),
					"existing variable %q must not be reused in short variable declaration; declare new variables separately and assign with =",
					identifier.Name,
				)
			}

			return true
		})
	}

	return nil, nil
}

func isTypeSwitchGuard(assignStatement *ast.AssignStmt) bool {
	if len(assignStatement.Lhs) != 1 || len(assignStatement.Rhs) != 1 {
		return false
	}

	typeAssertExpression, ok := assignStatement.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return false
	}

	return typeAssertExpression.Type == nil
}
