package style

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// errVarnameSettings 控制 errvarname 规则的行为
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

		// 对应 LHS 中 error 位置的标识符
		if errorIndex >= len(assignStatement.Lhs) {
			continue
		}

		identifier, ok := assignStatement.Lhs[errorIndex].(*ast.Ident)
		if !ok {
			continue
		}

		// _ 始终允许
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

// errorReturnIndex 返回函数调用中 error 返回值在结果列表中的索引
// 如果函数不返回 error 则返回 -1
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

// isErrorType 判断类型是否为 error 接口
func isErrorType(targetType types.Type) bool {
	if targetType == nil {
		return false
	}

	return targetType.String() == "error"
}
