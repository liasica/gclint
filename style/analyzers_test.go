package style

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrShortDeclarationAnalyzer(t *testing.T) {
	testDataDirectory := analysistest.TestData()

	analysistest.Run(t, testDataDirectory, newErrShortDeclarationAnalyzer(), "errshortdecl")
}

func TestNamedReturnExplicitValueAnalyzer(t *testing.T) {
	testDataDirectory := analysistest.TestData()

	analysistest.Run(t, testDataDirectory, newNamedReturnExplicitValueAnalyzer(), "namedreturn")
}

func TestJSONMapKeyAnalyzer(t *testing.T) {
	testDataDirectory := analysistest.TestData()

	analysistest.Run(t, testDataDirectory, newJSONMapKeyAnalyzer(), "jsonmapkey")
}

func TestLayerDependencyAnalyzer(t *testing.T) {
	testDataDirectory := analysistest.TestData()

	analysistest.Run(
		t,
		testDataDirectory,
		newLayerDependencyAnalyzer([]dependencyRule{
			{
				Source: "layerguard/repository",
				Forbidden: []string{
					"layerguard/service",
				},
			},
		}),
		"layerguard/repository",
	)
}

func TestSemanticVariableReuseAnalyzer(t *testing.T) {
	testDataDirectory := analysistest.TestData()

	analysistest.Run(t, testDataDirectory, newSemanticVariableReuseAnalyzer(), "semanticreuse")
}
