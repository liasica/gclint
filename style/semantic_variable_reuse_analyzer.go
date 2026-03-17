package style

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type semanticVariableState struct {
	baselineTokens       []string
	relatedSourceObjects map[types.Object]struct{}
	variableName         string
}

func newSemanticVariableReuseAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "varreuse",
		Doc:  "report reusing a semantic variable as a container for a different business object",
		Run:  runSemanticVariableReuseAnalyzer,
	}
}

func runSemanticVariableReuseAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			functionDeclaration, ok := node.(*ast.FuncDecl)
			if ok {
				checkSemanticVariableReuse(pass, functionDeclaration.Body)
				return false
			}

			functionLiteral, ok := node.(*ast.FuncLit)
			if ok {
				checkSemanticVariableReuse(pass, functionLiteral.Body)
				return false
			}

			return true
		})
	}

	return nil, nil
}

func checkSemanticVariableReuse(pass *analysis.Pass, functionBody *ast.BlockStmt) {
	if functionBody == nil {
		return
	}

	semanticStateByObject := make(map[types.Object]semanticVariableState)

	inspectWithoutNestedFunctions(functionBody, func(node ast.Node) bool {
		switch currentNode := node.(type) {
		case *ast.AssignStmt:
			recordAssignStatementSemantics(pass, semanticStateByObject, currentNode)
		case *ast.DeclStmt:
			recordDeclarationSemantics(pass, semanticStateByObject, currentNode)
		}

		return true
	})
}

func recordAssignStatementSemantics(
	pass *analysis.Pass,
	semanticStateByObject map[types.Object]semanticVariableState,
	assignStatement *ast.AssignStmt,
) {
	for leftHandSideIndex, leftHandSideExpression := range assignStatement.Lhs {
		identifier, ok := leftHandSideExpression.(*ast.Ident)
		if !ok || identifier.Name == "_" {
			continue
		}

		identifierObject := objectOfIdentifier(pass, identifier)
		if identifierObject == nil {
			continue
		}

		sourceExpression := assignmentSourceExpression(assignStatement, leftHandSideIndex)
		recordSemanticAssignment(
			pass,
			semanticStateByObject,
			identifierObject,
			identifier.Name,
			sourceExpression,
			identifier.Pos(),
		)
	}
}

func recordDeclarationSemantics(
	pass *analysis.Pass,
	semanticStateByObject map[types.Object]semanticVariableState,
	declarationStatement *ast.DeclStmt,
) {
	generalDeclaration, ok := declarationStatement.Decl.(*ast.GenDecl)
	if !ok || generalDeclaration.Tok != token.VAR {
		return
	}

	for _, declarationSpecification := range generalDeclaration.Specs {
		valueSpecification, ok := declarationSpecification.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for valueNameIndex, valueName := range valueSpecification.Names {
			if valueName.Name == "_" {
				continue
			}

			valueObject := pass.TypesInfo.Defs[valueName]
			if valueObject == nil {
				continue
			}

			sourceExpression := valueSpecSourceExpression(valueSpecification, valueNameIndex)
			recordSemanticAssignment(
				pass,
				semanticStateByObject,
				valueObject,
				valueName.Name,
				sourceExpression,
				valueName.Pos(),
			)
		}
	}
}

func recordSemanticAssignment(
	pass *analysis.Pass,
	semanticStateByObject map[types.Object]semanticVariableState,
	variableObject types.Object,
	variableName string,
	sourceExpression ast.Expr,
	position token.Pos,
) {
	variableTokens := semanticTokensFromName(variableName)
	if len(variableTokens) == 0 {
		return
	}

	sourceTokens := semanticTokensFromExpression(sourceExpression)
	if len(sourceTokens) == 0 {
		return
	}

	currentState, hasCurrentState := semanticStateByObject[variableObject]
	if !hasCurrentState {
		baselineTokens := intersectTokens(variableTokens, sourceTokens)
		if len(baselineTokens) == 0 {
			return
		}

		semanticStateByObject[variableObject] = semanticVariableState{
			baselineTokens:       baselineTokens,
			relatedSourceObjects: relatedVariableObjectsFromExpression(pass, sourceExpression),
			variableName:         variableName,
		}

		return
	}

	if sourceExpressionUsesOnlyRelatedVariables(pass, sourceExpression, currentState.relatedSourceObjects) {
		return
	}

	if hasTokenOverlap(currentState.baselineTokens, sourceTokens) {
		return
	}

	pass.Reportf(
		position,
		"variable %q was already established for %s and should not be reused with %s",
		currentState.variableName,
		formatTokens(currentState.baselineTokens),
		formatTokens(sourceTokens),
	)
}

func assignmentSourceExpression(assignStatement *ast.AssignStmt, leftHandSideIndex int) ast.Expr {
	if len(assignStatement.Rhs) == 0 {
		return nil
	}

	if len(assignStatement.Rhs) == 1 {
		return assignStatement.Rhs[0]
	}

	if leftHandSideIndex >= len(assignStatement.Rhs) {
		return nil
	}

	return assignStatement.Rhs[leftHandSideIndex]
}

func valueSpecSourceExpression(valueSpecification *ast.ValueSpec, valueNameIndex int) ast.Expr {
	if len(valueSpecification.Values) == 0 {
		return nil
	}

	if len(valueSpecification.Values) == 1 {
		return valueSpecification.Values[0]
	}

	if valueNameIndex >= len(valueSpecification.Values) {
		return nil
	}

	return valueSpecification.Values[valueNameIndex]
}
