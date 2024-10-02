package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// UserConfig represents the application's configuration settings.
type UserConfig struct {
	LogLevel    string      `yaml:"logLevel"`
	AnimeConfig AnimeConfig `yaml:"anime"`
}

// AnimeConfig contains anime specific configuration
type AnimeConfig struct {
	TitleLanguage string `yaml:"titleLanguage"`
	DisplayLayout string `yaml:"displayLayout"`
}

// DefaultConfig returns a UserConfig populated with default values.
func DefaultConfig() *UserConfig {
	return &UserConfig{
		LogLevel: "info",
		AnimeConfig: AnimeConfig{
			TitleLanguage: "english",
			DisplayLayout: "list",
		},
	}
}

// getConfigFilePath returns the path to the configuration file.
// It first checks if the HISAME_CONFIG_FILE environment variable is set.
// If set, it uses its value as the config file path.
// Otherwise, it defaults to the standard configuration directory.
func getConfigFilePath() (string, error) {
	// Check for the HISAME_CONFIG_FILE environment variable
	envPath := os.Getenv("HISAME_CONFIG_FILE")
	if envPath != "" {
		// Expand environment variables and user home directory (~) if present
		expandedPath, err := expandPath(envPath)
		if err != nil {
			logrus.Errorf("Error expanding HISAME_CONFIG_FILE env var file path %q: %v", envPath, err)
			return "", fmt.Errorf("failed to expand config file path '%s': %w", envPath, err)
		}
		return expandedPath, nil
	}

	// Fallback to the default config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Errorf("Error getting user config dir: %v", err)
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	configPath := filepath.Join(configDir, "hisame", "config.yaml")
	return configPath, nil
}

// expandPath expands environment variables and the tilde (~) in the given path.
func expandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to determine home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Expand any other environment variables present in the path
	path = os.ExpandEnv(path)

	return path, nil
}

// LoadConfig loads the configuration from the config file or returns default values.
func LoadConfig() (*UserConfig, error) {
	// Initialise config with default values
	cfg := DefaultConfig()

	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist.  Return defaults
			return cfg, nil
		}
		return nil, err
	}

	// Unmarshal YAML data into cfg
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
