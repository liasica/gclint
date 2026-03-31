package style

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gclint", newPlugin)
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
		newRedeclareAnalyzer(),
		newNamedReturnExplicitValueAnalyzer(),
		newJSONMapKeyAnalyzer(),
		newLayerDependencyAnalyzer(pluginInstance.settings.DependencyRules),
		newSemanticVariableReuseAnalyzer(),
		newFuncParamLinebreakAnalyzer(pluginInstance.settings.MaxInlineParams),
		newErrVarnameAnalyzer(resolveErrVarnameSettings(pluginInstance.settings.ErrVarname)),
	}, nil
}

func (plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
