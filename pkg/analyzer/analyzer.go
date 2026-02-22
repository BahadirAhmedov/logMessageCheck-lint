package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mycustomlinter",
	Doc:  "Checks weather the log message satisfies some rules",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {

		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		obj := pass.TypesInfo.Uses[selectorExpr.Sel]
		currentLogger := obj.Pkg().Path()
		availableLoggers := []string{"log/slog", "go.uber.org/zap"}

		if slices.Contains(availableLoggers, currentLogger) {
			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			supportedLogLevels := []string{
				"Info", "Debug", "Warn", "Error", "Fatal",
			}

			supportedLoggerNames := []string{
				"logger", "slog", "zapLogger", "log", "loggerInstance", "zapInstance", "myLogger", "slogLogger",
			}

			
			if slices.Contains(supportedLoggerNames, ident.Name) && slices.Contains(supportedLogLevels, selectorExpr.Sel.Name) {
				params := callExpr.Args
				for _, elements := range params {

					if be, ok := elements.(*ast.BinaryExpr); ok {
						if !IsSensitiveDataExpr(be) {
							pass.Reportf(node.Pos(), "log message must not contain potentially sensitive data")
							return true
						}
						continue
					}

					basicLit, ok := elements.(*ast.BasicLit)
					if !ok {
						return true
					}

					if basicLit.Kind != token.STRING {
						return true
					}

					message, err := strconv.Unquote(basicLit.Value)
					if err != nil {
						return true
					}

					if !IsEnglish(message) {
						pass.Reportf(node.Pos(), "log message must be only in english")
						return true
					}

					if !IsLower(message) {
						pass.Reportf(node.Pos(), "log message must start with a lowercase letter")
						return true
					}

					if !IsLetterOrNumber(message) {
						pass.Reportf(node.Pos(), "log message must not contain special symbols or emoji")
						return true
					}

					if !IsSensitiveData(message) {
						pass.Reportf(node.Pos(), "log message must not contain potentially sensitive data")
						return true
					}
				}
			}

		}
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

func IsLower(message string) bool {

	trimmed := strings.TrimLeftFunc(message, unicode.IsSpace)
	if trimmed == "" {
		return true
	}

	for _, r := range trimmed {
		if unicode.IsLetter(r) {
			if unicode.IsLower(r) {
				return true
			}
			return false
		}
	}

	return true
}

func IsEnglish(message string) bool {
	for _, r := range message {
		if unicode.In(r, unicode.Cyrillic) {
			return false
		}
	}
	return true
}

func IsLetterOrNumber(message string) bool {
	fmt.Println(message)
	for _, r := range message {
		if unicode.In(r, unicode.Latin) {
			continue
		}
		if r == ' ' {
			continue
		}
		if unicode.In(r, unicode.Number) {
			continue
		}
		return false
	}

	return true
}

func IsSensitiveData(message string) bool {
	if isSafeSensitivePhrase(message) {
		return true
	}

	if kw := isBareSensitiveKeyword(message); kw != "" {
		return false
	}

	if kw := leakingByLiteralPattern(message); kw != "" {
		return false
	}

	return true
}

func IsSensitiveDataExpr(expr ast.Expr) bool {
	found, _ := findSensitiveIdentInExpr(expr)
	return !found
}

var sensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token", "jwt", "bearer",
	"secret",
	"apikey", "api_key", "api-key",
	"private_key", "privatekey",
	"authorization", "cookie", "session",
	"access_key", "accesskey",
	"refresh_token", "refreshtoken",
}

var safePhrases = []string{
	"token validated",
	"token invalid",
	"token expired",
	"token missing",
	"token verification failed",
	"token verification succeeded",
}

func isSafeSensitivePhrase(message string) bool {
	m := normalizeSensitive(message)
	for _, p := range safePhrases {
		if m == normalizeSensitive(p) {
			return true
		}
	}
	return false
}

func isBareSensitiveKeyword(message string) string {
	m := normalizeSensitive(message)
	for _, kw := range sensitiveKeywords {
		if m == normalizeSensitive(kw) {
			return kw
		}
	}
	return ""
}

func leakingByLiteralPattern(message string) string {
	low := strings.ToLower(message)
	norm := normalizeSensitive(message)

	hasKV := strings.Contains(low, ":") || strings.Contains(low, "=")
	hasPrintf := strings.Contains(low, "%")

	if !(hasKV || hasPrintf) {
		return ""
	}

	for _, kw := range sensitiveKeywords {
		if strings.Contains(norm, normalizeSensitive(kw)) {
			return kw
		}
	}
	return ""
}

func normalizeSensitive(s string) string {
	s = strings.ToLower(s)
	return strings.NewReplacer("_", "", "-", "", " ", "", "\t", "", "\n", "", "\r", "").Replace(s)
}

func findSensitiveIdentInExpr(e ast.Expr) (bool, string) {
	found := false
	var name string

	ast.Inspect(e, func(n ast.Node) bool {
		if found {
			return false
		}
		switch x := n.(type) {
		case *ast.Ident:
			if matchesSensitiveName(x.Name) {
				found = true
				name = x.Name
				return false
			}
		case *ast.SelectorExpr:
			if matchesSensitiveName(x.Sel.Name) {
				found = true
				name = x.Sel.Name
				return false
			}
		}
		return true
	})

	return found, name
}

func matchesSensitiveName(name string) bool {
	n := normalizeSensitive(name)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(n, normalizeSensitive(kw)) {
			return true
		}
	}
	return false
}

func New(conf any) ([]*analysis.Analyzer, error) {
    return []*analysis.Analyzer{Analyzer}, nil
}