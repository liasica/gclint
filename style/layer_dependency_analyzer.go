package style

import (
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func newLayerDependencyAnalyzer(dependencyRules []dependencyRule) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "layerdep",
		Doc:  "report imports from higher-level packages into lower-level packages",
		Run: func(pass *analysis.Pass) (any, error) {
			return runLayerDependencyAnalyzer(pass, dependencyRules)
		},
	}
}

func runLayerDependencyAnalyzer(pass *analysis.Pass, dependencyRules []dependencyRule) (any, error) {
	if len(dependencyRules) == 0 {
		return nil, nil
	}

	currentPackagePath := pass.Pkg.Path()

	for _, dr := range dependencyRules {
		if !matchesPackagePrefix(currentPackagePath, dr.Source) {
			continue
		}

		for _, file := range pass.Files {
			for _, importSpec := range file.Imports {
				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					continue
				}

				for _, forbiddenPackage := range dr.Forbidden {
					if !matchesPackagePrefix(importPath, forbiddenPackage) {
						continue
					}

					pass.Reportf(
						importSpec.Path.Pos(),
						"package %q must not import higher-level package %q",
						currentPackagePath,
						importPath,
					)
				}
			}
		}
	}

	return nil, nil
}
