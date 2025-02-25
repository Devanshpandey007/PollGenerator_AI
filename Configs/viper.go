package Configs

import (
	"embed"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Embed all YAML files matching either config.*.yaml or provider.*.yaml,
//
//go:embed provider.*.yaml config.*.yaml
var embeddedFiles embed.FS

func NewViper(environment string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Files will be merged; order matters if keys conflict.
	// Ensure your YAML files use distinct top-level keys.
	filesToLoad := []string{
		"config.yaml",
		fmt.Sprintf("config.%s.yaml", environment),
		"provider.yaml",
		fmt.Sprintf("provider.%s.yaml", environment),
	}

	for _, f := range filesToLoad {
		data, err := embeddedFiles.ReadFile(f)
		if err != nil {
			// Skip missing files.
			continue
		}

		// Unmarshal YAML data into a temporary map.
		var fileMap map[string]interface{}
		if err := yaml.Unmarshal(data, &fileMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file %s: %w", f, err)
		}

		// Merge the map into viper.
		if err := v.MergeConfigMap(fileMap); err != nil {
			return nil, fmt.Errorf("failed to merge file %s: %w", f, err)
		}
	}

	// Automatically pick up environment variables,
	// replacing dots with underscores.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
