package style

import (
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var genericSemanticTokens = map[string]struct{}{
	"append":    {},
	"build":     {},
	"bytes":     {},
	"byte":      {},
	"body":      {},
	"content":   {},
	"create":    {},
	"current":   {},
	"ctx":       {},
	"data":      {},
	"decode":    {},
	"detail":    {},
	"details":   {},
	"encode":    {},
	"err":       {},
	"fetch":     {},
	"get":       {},
	"info":      {},
	"item":      {},
	"items":     {},
	"json":      {},
	"list":      {},
	"load":      {},
	"log":       {},
	"logs":      {},
	"marshal":   {},
	"message":   {},
	"messages":  {},
	"msg":       {},
	"name":      {},
	"names":     {},
	"normalize": {},
	"ok":        {},
	"parse":     {},
	"payload":   {},
	"raw":       {},
	"read":      {},
	"record":    {},
	"records":   {},
	"request":   {},
	"requests":  {},
	"req":       {},
	"response":  {},
	"responses": {},
	"resp":      {},
	"result":    {},
	"results":   {},
	"root":      {},
	"rune":      {},
	"selector":  {},
	"set":       {},
	"tmp":       {},
	"temp":      {},
	"token":     {},
	"tokens":    {},
	"unmarshal": {},
	"update":    {},
	"value":     {},
	"values":    {},
	"write":     {},
}

func inspectWithoutNestedFunctions(root ast.Node, visit func(ast.Node) bool) {
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		if node != root {
			if _, ok := node.(*ast.FuncLit); ok {
				return false
			}
		}

		return visit(node)
	})
}

func objectOfIdentifier(pass *analysis.Pass, identifier *ast.Ident) types.Object {
	if identifierObject := pass.TypesInfo.Defs[identifier]; identifierObject != nil {
		return identifierObject
	}

	return pass.TypesInfo.Uses[identifier]
}

func stringConstantValue(pass *analysis.Pass, expression ast.Expr) (string, bool) {
	typeAndValue, ok := pass.TypesInfo.Types[expression]
	if !ok || typeAndValue.Value == nil {
		return "", false
	}

	if typeAndValue.Value.Kind() != constant.String {
		return "", false
	}

	return constant.StringVal(typeAndValue.Value), true
}

func containsChinese(text string) bool {
	for _, currentRune := range text {
		if unicode.Is(unicode.Han, currentRune) {
			return true
		}
	}

	return false
}

func matchesPackagePrefix(packagePath string, prefix string) bool {
	return packagePath == prefix || strings.HasPrefix(packagePath, prefix+"/")
}

func semanticTokensFromName(name string) []string {
	lowercaseTokens := splitIdentifierTokens(name)
	semanticTokens := make([]string, 0, len(lowercaseTokens))

	for _, lowercaseToken := range lowercaseTokens {
		if len(lowercaseToken) <= 1 {
			continue
		}

		if _, ok := genericSemanticTokens[lowercaseToken]; ok {
			continue
		}

		semanticTokens = append(semanticTokens, lowercaseToken)
	}

	return uniqueTokens(semanticTokens)
}

func semanticTokensFromExpression(expression ast.Expr) []string {
	if expression == nil {
		return nil
	}

	switch currentExpression := expression.(type) {
	case *ast.CallExpr:
		return semanticTokensFromExpression(currentExpression.Fun)
	case *ast.SelectorExpr:
		selectorTokens := semanticTokensFromExpression(currentExpression.X)
		selectorTokens = append(selectorTokens, semanticTokensFromName(currentExpression.Sel.Name)...)
		return uniqueTokens(selectorTokens)
	case *ast.Ident:
		return semanticTokensFromName(currentExpression.Name)
	case *ast.CompositeLit:
		return semanticTokensFromExpression(currentExpression.Type)
	case *ast.StarExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.UnaryExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.IndexExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.IndexListExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.SliceExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.ParenExpr:
		return semanticTokensFromExpression(currentExpression.X)
	case *ast.ArrayType:
		return semanticTokensFromExpression(currentExpression.Elt)
	case *ast.MapType:
		return semanticTokensFromExpression(currentExpression.Value)
	case *ast.StructType:
		return nil
	case *ast.InterfaceType:
		return nil
	}

	return nil
}

func splitIdentifierTokens(name string) []string {
	if name == "" {
		return nil
	}

	tokens := make([]string, 0)
	currentToken := make([]rune, 0, len(name))
	previousRuneWasLower := false

	for _, currentRune := range name {
		if currentRune == '_' || currentRune == '-' || unicode.IsDigit(currentRune) {
			if len(currentToken) > 0 {
				tokens = append(tokens, strings.ToLower(string(currentToken)))
				currentToken = currentToken[:0]
			}

			previousRuneWasLower = false
			continue
		}

		if unicode.IsUpper(currentRune) && previousRuneWasLower && len(currentToken) > 0 {
			tokens = append(tokens, strings.ToLower(string(currentToken)))
			currentToken = currentToken[:0]
		}

		currentToken = append(currentToken, currentRune)
		previousRuneWasLower = unicode.IsLower(currentRune)
	}

	if len(currentToken) > 0 {
		tokens = append(tokens, strings.ToLower(string(currentToken)))
	}

	return tokens
}

func uniqueTokens(tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}

	seenTokens := make(map[string]struct{}, len(tokens))
	uniqueTokenList := make([]string, 0, len(tokens))

	for _, token := range tokens {
		if _, ok := seenTokens[token]; ok {
			continue
		}

		seenTokens[token] = struct{}{}
		uniqueTokenList = append(uniqueTokenList, token)
	}

	return uniqueTokenList
}

func intersectTokens(leftTokens []string, rightTokens []string) []string {
	if len(leftTokens) == 0 || len(rightTokens) == 0 {
		return nil
	}

	rightTokenSet := make(map[string]struct{}, len(rightTokens))
	for _, rightToken := range rightTokens {
		rightTokenSet[rightToken] = struct{}{}
	}

	intersection := make([]string, 0)
	for _, leftToken := range leftTokens {
		if _, ok := rightTokenSet[leftToken]; ok {
			intersection = append(intersection, leftToken)
		}
	}

	return uniqueTokens(intersection)
}

func hasTokenOverlap(leftTokens []string, rightTokens []string) bool {
	if len(leftTokens) == 0 || len(rightTokens) == 0 {
		return false
	}

	rightTokenSet := make(map[string]struct{}, len(rightTokens))
	for _, rightToken := range rightTokens {
		rightTokenSet[rightToken] = struct{}{}
	}

	for _, leftToken := range leftTokens {
		if _, ok := rightTokenSet[leftToken]; ok {
			return true
		}
	}

	return false
}

func formatTokens(tokens []string) string {
	if len(tokens) == 0 {
		return "[]"
	}

	sortedTokens := append([]string(nil), tokens...)
	sort.Strings(sortedTokens)

	return "[" + strings.Join(sortedTokens, ", ") + "]"
}
