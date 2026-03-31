package style

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// defaultMaxInlineParams 参数个数小于此值时禁止换行
const defaultMaxInlineParams = 5

func newFuncParamLinebreakAnalyzer(maxInlineParams int) *analysis.Analyzer {
	if maxInlineParams <= 0 {
		maxInlineParams = defaultMaxInlineParams
	}

	return &analysis.Analyzer{
		Name: "funcparamlinebreak",
		Doc:  "forbid line breaks in function parameter lists when the parameter count is below a configurable threshold",
		Run:  makeFuncParamLinebreakRunner(maxInlineParams),
	}
}

func makeFuncParamLinebreakRunner(maxInlineParams int) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch currentNode := node.(type) {
				case *ast.FuncDecl:
					checkFuncParamLinebreak(pass, currentNode.Type, maxInlineParams)
				case *ast.FuncLit:
					checkFuncParamLinebreak(pass, currentNode.Type, maxInlineParams)
				}

				return true
			})
		}

		return nil, nil
	}
}

func checkFuncParamLinebreak(pass *analysis.Pass, funcType *ast.FuncType, maxInlineParams int) {
	if funcType == nil || funcType.Params == nil {
		return
	}

	paramCount := countFuncParams(funcType.Params)
	if paramCount == 0 || paramCount >= maxInlineParams {
		return
	}

	openParenLine := pass.Fset.Position(funcType.Params.Opening).Line
	closeParenLine := pass.Fset.Position(funcType.Params.Closing).Line

	if openParenLine != closeParenLine {
		pass.Reportf(
			funcType.Params.Opening,
			"function has %d parameters (threshold %d), parameter list must not span multiple lines",
			paramCount,
			maxInlineParams,
		)
	}
}

// countFuncParams 统计函数参数列表中的参数个数
func countFuncParams(fieldList *ast.FieldList) int {
	if fieldList == nil {
		return 0
	}

	count := 0

	for _, field := range fieldList.List {
		if len(field.Names) == 0 {
			// 匿名参数（例如接口方法中的 string, error）
			count++
		} else {
			count += len(field.Names)
		}
	}

	return count
}
