package style

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type namedReturnVariable struct {
	object types.Object
}

func newNamedReturnExplicitValueAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "namedreturn",
		Doc:  "report explicit return values after named return variables have already been assigned",
		Run:  runNamedReturnExplicitValueAnalyzer,
	}
}

func runNamedReturnExplicitValueAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			functionDeclaration, isFunctionDeclaration := node.(*ast.FuncDecl)
			if isFunctionDeclaration {
				checkFunctionBody(pass, functionDeclaration.Type, functionDeclaration.Body)
				return false
			}

			functionLiteral, isFunctionLiteral := node.(*ast.FuncLit)
			if isFunctionLiteral {
				checkFunctionBody(pass, functionLiteral.Type, functionLiteral.Body)
				return false
			}

			return true
		})
	}

	return nil, nil
}

func checkFunctionBody(pass *analysis.Pass, functionType *ast.FuncType, functionBody *ast.BlockStmt) {
	if functionType == nil || functionType.Results == nil || functionBody == nil {
		return
	}

	namedReturnVariables := collectNamedReturnVariables(pass, functionType.Results)
	if len(namedReturnVariables) == 0 {
		return
	}

	namedReturnObjects := make(map[types.Object]struct{}, len(namedReturnVariables))
	for _, nrv := range namedReturnVariables {
		namedReturnObjects[nrv.object] = struct{}{}
	}

	assignedNamedReturnPositions := collectAssignedNamedReturnPositions(pass, functionBody, namedReturnObjects)

	inspectWithoutNestedFunctions(functionBody, func(node ast.Node) bool {
		returnStatement, ok := node.(*ast.ReturnStmt)
		if !ok || len(returnStatement.Results) == 0 {
			return true
		}

		if !hasAssignedNamedReturnBeforePosition(assignedNamedReturnPositions, returnStatement.Pos()) {
			return true
		}

		pass.Reportf(
			returnStatement.Return,
			"named return values were assigned before this explicit return; use bare return instead",
		)

		return true
	})
}

func collectNamedReturnVariables(pass *analysis.Pass, resultFields *ast.FieldList) []namedReturnVariable {
	namedReturnVariables := make([]namedReturnVariable, 0)

	for _, resultField := range resultFields.List {
		for _, resultName := range resultField.Names {
			if resultName.Name == "_" {
				continue
			}

			resultObject := pass.TypesInfo.Defs[resultName]
			if resultObject == nil {
				continue
			}

			namedReturnVariables = append(namedReturnVariables, namedReturnVariable{
				object: resultObject,
			})
		}
	}

	return namedReturnVariables
}

func collectAssignedNamedReturnPositions(
	pass *analysis.Pass,
	functionBody *ast.BlockStmt,
	namedReturnObjects map[types.Object]struct{},
) map[types.Object]token.Pos {
	assignedNamedReturnPositions := make(map[types.Object]token.Pos, len(namedReturnObjects))

	inspectWithoutNestedFunctions(functionBody, func(node ast.Node) bool {
		assignStatement, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for _, leftHandSideExpression := range assignStatement.Lhs {
			identifier, ok := leftHandSideExpression.(*ast.Ident)
			if !ok {
				continue
			}

			identifierObject := objectOfIdentifier(pass, identifier)
			if identifierObject == nil {
				continue
			}

			if _, ok = namedReturnObjects[identifierObject]; !ok {
				continue
			}

			existingPosition, exists := assignedNamedReturnPositions[identifierObject]
			if !exists || identifier.Pos() < existingPosition {
				assignedNamedReturnPositions[identifierObject] = identifier.Pos()
			}
		}

		return true
	})

	return assignedNamedReturnPositions
}

func hasAssignedNamedReturnBeforePosition(
	assignedNamedReturnPositions map[types.Object]token.Pos,
	position token.Pos,
) bool {
	for _, assignedPosition := range assignedNamedReturnPositions {
		if assignedPosition < position {
			return true
		}
	}

	return false
}
