package Configs

import (
	"bytes"
	"embed"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Embed all YAML files matching either config.*.yaml or provider.*.yaml,
//
//go:embed provider.*.yaml config.*.yaml
var embeddedFiles embed.FS

func NewViper(environment string) (*viper.Viper, error) {

	v := viper.New()
	v.SetConfigType("yaml")

	// We'll load multiple files, in a specified order, merging them as we go.
	// 1) provider.<env>.yaml
	filesToLoad := []string{
		"config.yaml",
		fmt.Sprintf("config.%s.yaml", environment),
		"provider.yaml",
		fmt.Sprintf("provider.%s.yaml", environment),
	}

	for _, f := range filesToLoad {
		if data, err := embeddedFiles.ReadFile(f); err == nil {
			// If the file exists in the embedded FS, merge it
			if v.ConfigFileUsed() == "" {
				// No config loaded yet, so read the config fresh
				if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
					return nil, fmt.Errorf("failed to load file %s: %w", f, err)
				}
			} else {
				// Already have a config; merge additional files
				if err := v.MergeConfig(bytes.NewReader(data)); err != nil {
					return nil, fmt.Errorf("failed to merge file %s: %w", f, err)
				}
			}
		}
		// If a file is missing, we just skip.
	}

	// Automatically pick up environment variables,
	// with dot replaced by underscore in the env var name
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
