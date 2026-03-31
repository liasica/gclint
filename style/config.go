package style

import (
	"strings"

	"github.com/golangci/plugin-module-register/register"
)

type pluginSettings struct {
	DependencyRules []dependencyRule  `json:"dependency_rules"`
	MaxInlineParams int               `json:"max_inline_params"`
	ErrVarname      *errVarnameConfig `json:"err_varname"`
}

type errVarnameConfig struct {
	Enabled      *bool    `json:"enabled"`
	AllowedNames []string `json:"allowed_names"`
}

type dependencyRule struct {
	Source    string   `json:"source"`
	Forbidden []string `json:"forbidden"`
}

type rawPluginConfig struct {
	Type            string            `json:"type"`
	Description     string            `json:"description"`
	OriginalURL     string            `json:"original-url"`
	Path            string            `json:"path"`
	Settings        pluginSettings    `json:"settings"`
	DependencyRules []dependencyRule  `json:"dependency_rules"`
	MaxInlineParams int               `json:"max_inline_params"`
	ErrVarname      *errVarnameConfig `json:"err_varname"`
}

func decodePluginSettings(rawSettings any) (pluginSettings, error) {
	if rawSettings == nil {
		return pluginSettings{}, nil
	}

	rawConfig, rawConfigError := register.DecodeSettings[rawPluginConfig](rawSettings)
	if rawConfigError == nil {
		decodedSettings := rawConfig.Settings
		if len(decodedSettings.DependencyRules) == 0 {
			decodedSettings.DependencyRules = rawConfig.DependencyRules
		}

		if decodedSettings.MaxInlineParams == 0 {
			decodedSettings.MaxInlineParams = rawConfig.MaxInlineParams
		}

		if decodedSettings.ErrVarname == nil {
			decodedSettings.ErrVarname = rawConfig.ErrVarname
		}

		return normalizePluginSettings(decodedSettings), nil
	}

	decodedSettings, settingsError := register.DecodeSettings[pluginSettings](rawSettings)
	if settingsError != nil {
		return pluginSettings{}, rawConfigError
	}

	return normalizePluginSettings(decodedSettings), nil
}

func normalizePluginSettings(settings pluginSettings) pluginSettings {
	if settings.MaxInlineParams <= 0 {
		settings.MaxInlineParams = defaultMaxInlineParams
	}

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

// resolveErrVarnameSettings 将配置转换为 analyzer 可用的 settings
func resolveErrVarnameSettings(config *errVarnameConfig) errVarnameSettings {
	defaults := defaultErrVarnameSettings()

	if config == nil {
		return defaults
	}

	if config.Enabled != nil {
		defaults.Enabled = *config.Enabled
	}

	if len(config.AllowedNames) > 0 {
		allowedNames := make(map[string]struct{}, len(config.AllowedNames))
		for _, name := range config.AllowedNames {
			trimmed := strings.TrimSpace(name)
			if trimmed != "" {
				allowedNames[trimmed] = struct{}{}
			}
		}

		if len(allowedNames) > 0 {
			defaults.AllowedNames = allowedNames
		}
	}

	return defaults
}
