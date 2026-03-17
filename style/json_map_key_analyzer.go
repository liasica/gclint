package style

import (
	"go/ast"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func newJSONMapKeyAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "chinesekey",
		Doc:  "report Chinese JSON tag keys and Chinese string keys in map literals",
		Run:  runJSONMapKeyAnalyzer,
	}
}

func runJSONMapKeyAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			structField, ok := node.(*ast.Field)
			if ok {
				checkJSONTag(pass, structField)
			}

			compositeLiteral, ok := node.(*ast.CompositeLit)
			if ok {
				checkMapLiteral(pass, compositeLiteral)
			}

			return true
		})
	}

	return nil, nil
}

func checkJSONTag(pass *analysis.Pass, structField *ast.Field) {
	if structField.Tag == nil {
		return
	}

	unquotedTag, err := strconv.Unquote(structField.Tag.Value)
	if err != nil {
		return
	}

	jsonTagValue := reflect.StructTag(unquotedTag).Get("json")
	if jsonTagValue == "" || jsonTagValue == "-" {
		return
	}

	jsonTagKey := strings.Split(jsonTagValue, ",")[0]
	if jsonTagKey == "" || !containsChinese(jsonTagKey) {
		return
	}

	pass.Reportf(structField.Tag.Pos(), "JSON tag key must not contain Chinese: %q", jsonTagKey)
}

func checkMapLiteral(pass *analysis.Pass, compositeLiteral *ast.CompositeLit) {
	compositeLiteralType := pass.TypesInfo.TypeOf(compositeLiteral)
	if compositeLiteralType == nil {
		return
	}

	mapType, ok := compositeLiteralType.Underlying().(*types.Map)
	if !ok {
		return
	}

	if mapType.Key().String() != "string" {
		return
	}

	for _, compositeElement := range compositeLiteral.Elts {
		keyValueExpression, ok := compositeElement.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		mapKeyValue, ok := stringConstantValue(pass, keyValueExpression.Key)
		if !ok || !containsChinese(mapKeyValue) {
			continue
		}

		pass.Reportf(keyValueExpression.Key.Pos(), "map key must not contain Chinese: %q", mapKeyValue)
	}
}
