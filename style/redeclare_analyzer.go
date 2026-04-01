package style

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func newRedeclareAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "redeclare",
		Doc:  "report short variable declarations that reuse an existing variable from the current function scope chain",
		Run:  runRedeclareAnalyzer,
	}
}

func runRedeclareAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch currentNode := node.(type) {
			case *ast.FuncDecl:
				runRedeclareAnalyzerInFunction(pass, currentNode.Type, currentNode.Body)
			case *ast.FuncLit:
				runRedeclareAnalyzerInFunction(pass, currentNode.Type, currentNode.Body)
			}

			return true
		})
	}

	return nil, nil
}

func runRedeclareAnalyzerInFunction(pass *analysis.Pass, functionType *ast.FuncType, functionBody *ast.BlockStmt) {
	if functionType == nil || functionBody == nil {
		return
	}

	functionScopes := collectFunctionScopes(pass, functionType, functionBody)

	inspectWithoutNestedFunctions(functionBody, func(node ast.Node) bool {
		switch currentNode := node.(type) {
		case *ast.AssignStmt:
			if currentNode.Tok != token.DEFINE {
				return true
			}

			if isTypeSwitchGuard(currentNode) {
				return true
			}

			reportRedeclaredIdentifiers(pass, currentNode.Lhs, functionScopes)
		case *ast.RangeStmt:
			if currentNode.Tok != token.DEFINE {
				return true
			}

			leftHandSideExpressions := make([]ast.Expr, 0, 2)

			if currentNode.Key != nil {
				leftHandSideExpressions = append(leftHandSideExpressions, currentNode.Key)
			}

			if currentNode.Value != nil {
				leftHandSideExpressions = append(leftHandSideExpressions, currentNode.Value)
			}

			reportRedeclaredIdentifiers(pass, leftHandSideExpressions, functionScopes)
		}

		return true
	})
}

func collectFunctionScopes(pass *analysis.Pass, functionType *ast.FuncType, functionBody *ast.BlockStmt) map[*types.Scope]struct{} {
	functionScopes := make(map[*types.Scope]struct{})

	addFunctionScope(functionScopes, pass.TypesInfo.Scopes[functionType])
	inspectWithoutNestedFunctions(functionBody, func(node ast.Node) bool {
		addFunctionScope(functionScopes, pass.TypesInfo.Scopes[node])

		return true
	})

	return functionScopes
}

func addFunctionScope(functionScopes map[*types.Scope]struct{}, scope *types.Scope) {
	if scope == nil {
		return
	}

	functionScopes[scope] = struct{}{}
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

func reportRedeclaredIdentifiers(pass *analysis.Pass, leftHandSideExpressions []ast.Expr, functionScopes map[*types.Scope]struct{}) {
	for _, leftHandSideExpression := range leftHandSideExpressions {
		identifier, ok := leftHandSideExpression.(*ast.Ident)
		if !ok || identifier.Name == "_" {
			continue
		}

		reusedVariable := reusedVariableInFunctionScopeChain(pass, identifier, functionScopes)
		if reusedVariable == nil {
			continue
		}

		pass.Reportf(
			identifier.Pos(),
			"existing variable %q must not be reused in short variable declaration; use a distinct name or assign to the existing variable with =",
			identifier.Name,
		)
	}
}

func reusedVariableInFunctionScopeChain(pass *analysis.Pass, identifier *ast.Ident, functionScopes map[*types.Scope]struct{}) types.Object {
	identifierObject := objectOfIdentifier(pass, identifier)
	if identifierObject == nil {
		return nil
	}

	if _, ok := identifierObject.(*types.Var); !ok {
		return nil
	}

	if pass.TypesInfo.Defs[identifier] == nil {
		if !scopeBelongsToFunction(functionScopes, identifierObject.Parent()) {
			return nil
		}

		return identifierObject
	}

	definedObject := pass.TypesInfo.Defs[identifier]

	currentScope := definedObject.Parent()
	if currentScope == nil {
		return nil
	}

	for currentScope = currentScope.Parent(); currentScope != nil; currentScope = currentScope.Parent() {
		if !scopeBelongsToFunction(functionScopes, currentScope) {
			break
		}

		outerObject := currentScope.Lookup(identifier.Name)
		if outerObject != nil && outerObject.Pos() < identifier.Pos() {
			if _, ok := outerObject.(*types.Var); ok {
				return outerObject
			}
		}
	}

	return nil
}

func scopeBelongsToFunction(functionScopes map[*types.Scope]struct{}, scope *types.Scope) bool {
	if scope == nil {
		return false
	}

	_, ok := functionScopes[scope]

	return ok
}
