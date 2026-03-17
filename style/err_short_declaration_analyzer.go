package style

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func newErrShortDeclarationAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "errshort",
		Doc:  "report short variable declarations that reuse an existing err in the same scope",
		Run:  runErrShortDeclarationAnalyzer,
	}
}

func runErrShortDeclarationAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			assignStatement, ok := node.(*ast.AssignStmt)
			if !ok || assignStatement.Tok != token.DEFINE {
				return true
			}

			for _, leftHandSideExpression := range assignStatement.Lhs {
				var identifier *ast.Ident
				identifier, ok = leftHandSideExpression.(*ast.Ident)
				if !ok || identifier.Name != "err" {
					continue
				}

				if pass.TypesInfo.Defs[identifier] != nil {
					continue
				}

				pass.Reportf(
					identifier.Pos(),
					"existing err must not be reused in short variable declaration; declare new variables separately and assign to err with =",
				)

				return true
			}

			return true
		})
	}

	return nil, nil
}
