package style

import (
	"strings"

	"github.com/golangci/plugin-module-register/register"
)

type pluginSettings struct {
	DependencyRules []dependencyRule `json:"dependency_rules"`
}

type dependencyRule struct {
	Source    string   `json:"source"`
	Forbidden []string `json:"forbidden"`
}

type rawPluginConfig struct {
	Type            string           `json:"type"`
	Description     string           `json:"description"`
	OriginalURL     string           `json:"original-url"`
	Path            string           `json:"path"`
	Settings        pluginSettings   `json:"settings"`
	DependencyRules []dependencyRule `json:"dependency_rules"`
}

func decodePluginSettings(rawSettings any) (pluginSettings, error) {
	if rawSettings == nil {
		return pluginSettings{}, nil
	}

	rawConfig, rawConfigError := register.DecodeSettings[rawPluginConfig](rawSettings)
	if rawConfigError == nil {
		settings := rawConfig.Settings
		if len(settings.DependencyRules) == 0 {
			settings.DependencyRules = rawConfig.DependencyRules
		}

		return normalizePluginSettings(settings), nil
	}

	settings, settingsError := register.DecodeSettings[pluginSettings](rawSettings)
	if settingsError != nil {
		return pluginSettings{}, rawConfigError
	}

	return normalizePluginSettings(settings), nil
}

func normalizePluginSettings(settings pluginSettings) pluginSettings {
	normalizedDependencyRules := make([]dependencyRule, 0, len(settings.DependencyRules))

	for _, currentDependencyRule := range settings.DependencyRules {
		trimmedSource := strings.TrimSpace(currentDependencyRule.Source)
		if trimmedSource == "" {
			continue
		}

		normalizedForbiddenPackages := make([]string, 0, len(currentDependencyRule.Forbidden))
		for _, forbiddenPackage := range currentDependencyRule.Forbidden {
			trimmedForbiddenPackage := strings.TrimSpace(forbiddenPackage)
			if trimmedForbiddenPackage == "" {
				continue
			}

			normalizedForbiddenPackages = append(normalizedForbiddenPackages, trimmedForbiddenPackage)
		}

		if len(normalizedForbiddenPackages) == 0 {
			continue
		}

		normalizedDependencyRules = append(normalizedDependencyRules, dependencyRule{
			Source:    trimmedSource,
			Forbidden: normalizedForbiddenPackages,
		})
	}

	settings.DependencyRules = normalizedDependencyRules

	return settings
}
