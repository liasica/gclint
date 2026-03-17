package style

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("style", newPlugin)
}

type plugin struct {
	settings pluginSettings
}

func newPlugin(rawSettings any) (register.LinterPlugin, error) {
	settings, err := decodePluginSettings(rawSettings)
	if err != nil {
		return nil, err
	}

	return &plugin{
		settings: settings,
	}, nil
}

func (pluginInstance plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		newErrShortDeclarationAnalyzer(),
		newRedeclareAnalyzer(),
		newNamedReturnExplicitValueAnalyzer(),
		newJSONMapKeyAnalyzer(),
		newLayerDependencyAnalyzer(pluginInstance.settings.DependencyRules),
		newSemanticVariableReuseAnalyzer(),
	}, nil
}

func (plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
