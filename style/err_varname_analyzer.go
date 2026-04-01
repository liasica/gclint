package style

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// errVarnameSettings controls the behavior of the errvarname rule
type errVarnameSettings struct {
	Enabled      bool
	AllowedNames map[string]struct{}
}

func defaultErrVarnameSettings() errVarnameSettings {
	return errVarnameSettings{
		Enabled: true,
		AllowedNames: map[string]struct{}{
			"err": {},
		},
	}
}

func newErrVarnameAnalyzer(settings errVarnameSettings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "errvarname",
		Doc:  "forbid using non-err variable names to receive error returns from function calls",
		Run:  makeErrVarnameRunner(settings),
	}
}

func makeErrVarnameRunner(settings errVarnameSettings) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		if !settings.Enabled {
			return nil, nil
		}

		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				assignStatement, ok := node.(*ast.AssignStmt)
				if !ok {
					return true
				}

				if assignStatement.Tok != token.DEFINE && assignStatement.Tok != token.ASSIGN {
					return true
				}

				checkErrVarnameInAssignment(pass, assignStatement, settings.AllowedNames)

				return true
			})
		}

		return nil, nil
	}
}

func checkErrVarnameInAssignment(pass *analysis.Pass, assignStatement *ast.AssignStmt, allowedNames map[string]struct{}) {
	for _, rightHandSideExpression := range assignStatement.Rhs {
		callExpression, ok := rightHandSideExpression.(*ast.CallExpr)
		if !ok {
			continue
		}

		errorIndex := errorReturnIndex(pass, callExpression)
		if errorIndex < 0 {
			continue
		}

		// identifier at the error position in LHS
		if errorIndex >= len(assignStatement.Lhs) {
			continue
		}

		var identifier *ast.Ident

		identifier, ok = assignStatement.Lhs[errorIndex].(*ast.Ident)
		if !ok {
			continue
		}

		// _ is always allowed
		if identifier.Name == "_" {
			continue
		}

		if _, allowed := allowedNames[identifier.Name]; allowed {
			continue
		}

		allowedList := make([]string, 0, len(allowedNames))
		for name := range allowedNames {
			allowedList = append(allowedList, name)
		}

		pass.Reportf(
			identifier.Pos(),
			"error return must be received by %q, not %q",
			strings.Join(allowedList, "\" or \""),
			identifier.Name,
		)
	}
}

// errorReturnIndex returns the index of the error return value in the result list of a function call
// returns -1 if the function does not return an error
func errorReturnIndex(pass *analysis.Pass, callExpression *ast.CallExpr) int {
	callType := pass.TypesInfo.TypeOf(callExpression.Fun)
	if callType == nil {
		return -1
	}

	signature, ok := callType.Underlying().(*types.Signature)
	if !ok {
		return -1
	}

	results := signature.Results()
	if results == nil || results.Len() == 0 {
		return -1
	}

	lastIndex := results.Len() - 1
	lastResult := results.At(lastIndex)

	if !isErrorType(lastResult.Type()) {
		return -1
	}

	return lastIndex
}

// isErrorType checks whether the type is the error interface
func isErrorType(targetType types.Type) bool {
	if targetType == nil {
		return false
	}

	return targetType.String() == "error"
}
