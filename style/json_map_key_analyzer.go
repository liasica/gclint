package style

import (
	"encoding/json"
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
		Doc:  "report Chinese keys in JSON tags, persistent maps, and raw JSON string constants",
		Run:  runJSONMapKeyAnalyzer,
	}
}

func runJSONMapKeyAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			structField, isStructField := node.(*ast.Field)
			if isStructField {
				checkJSONTag(pass, structField)
			}

			assignStatement, isAssignStatement := node.(*ast.AssignStmt)
			if isAssignStatement {
				checkMapAssignment(pass, assignStatement)
			}

			compositeLiteral, isCompositeLiteral := node.(*ast.CompositeLit)
			if isCompositeLiteral {
				checkMapLiteral(pass, compositeLiteral)
			}

			if checkJSONStringConstant(pass, node) {
				return true
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
		keyValueExpression, isKeyValueExpression := compositeElement.(*ast.KeyValueExpr)
		if !isKeyValueExpression {
			continue
		}

		mapKeyValue, hasConstantValue := stringConstantValue(pass, keyValueExpression.Key)
		if !hasConstantValue || !containsChinese(mapKeyValue) {
			continue
		}

		pass.Reportf(keyValueExpression.Key.Pos(), "map key must not contain Chinese: %q", mapKeyValue)
	}
}

func checkMapAssignment(pass *analysis.Pass, assignStatement *ast.AssignStmt) {
	for _, leftHandSideExpression := range assignStatement.Lhs {
		indexExpression, isIndexExpression := leftHandSideExpression.(*ast.IndexExpr)
		if !isIndexExpression {
			continue
		}

		indexTargetType := pass.TypesInfo.TypeOf(indexExpression.X)
		if indexTargetType == nil {
			continue
		}

		mapType, isMapType := indexTargetType.Underlying().(*types.Map)
		if !isMapType || mapType.Key().String() != "string" {
			continue
		}

		mapKeyValue, hasConstantValue := stringConstantValue(pass, indexExpression.Index)
		if !hasConstantValue || !containsChinese(mapKeyValue) {
			continue
		}

		pass.Reportf(indexExpression.Index.Pos(), "map key must not contain Chinese: %q", mapKeyValue)
	}
}

func checkJSONStringConstant(pass *analysis.Pass, node ast.Node) bool {
	expression, isBasicLiteral := node.(*ast.BasicLit)
	if !isBasicLiteral {
		return false
	}

	jsonSource, hasConstantValue := stringConstantValue(pass, expression)
	if !hasConstantValue {
		return false
	}

	for _, jsonKey := range chineseJSONKeys(jsonSource) {
		pass.Reportf(expression.Pos(), "JSON string key must not contain Chinese: %q", jsonKey)
	}

	return false
}

func chineseJSONKeys(rawJSON string) []string {
	trimmedJSON := strings.TrimSpace(rawJSON)
	if trimmedJSON == "" {
		return nil
	}

	firstCharacter := trimmedJSON[0]
	if firstCharacter != '{' && firstCharacter != '[' {
		return nil
	}

	var decodedValue any
	if err := json.Unmarshal([]byte(trimmedJSON), &decodedValue); err != nil {
		return nil
	}

	collectedKeys := make([]string, 0)
	collectChineseJSONKeys(decodedValue, &collectedKeys)

	return uniqueTokens(collectedKeys)
}

func collectChineseJSONKeys(value any, collectedKeys *[]string) {
	switch decodedValue := value.(type) {
	case map[string]any:
		for key, nestedValue := range decodedValue {
			if containsChinese(key) {
				*collectedKeys = append(*collectedKeys, key)
			}

			collectChineseJSONKeys(nestedValue, collectedKeys)
		}
	case []any:
		for _, nestedValue := range decodedValue {
			collectChineseJSONKeys(nestedValue, collectedKeys)
		}
	}
}
