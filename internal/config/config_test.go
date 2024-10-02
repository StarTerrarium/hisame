package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_DefaultsWhenFileMissing(t *testing.T) {
	// Set up a temporary directory for config files
	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	// Remove any existing config file
	configPath, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("Failed to get config file path: %v", err)
	}
	os.Remove(configPath)

	// Test loading config when file doesn't exist
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check default values
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default LogLevel 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.AnimeConfig.TitleLanguage != "english" {
		t.Errorf("Expected default Anime TitleLanguage 'english', got '%s'", cfg.AnimeConfig.TitleLanguage)
	}
	if cfg.AnimeConfig.DisplayLayout != "list" {
		t.Errorf("Expected default Anime DisplayLayout 'list', got '%s'", cfg.AnimeConfig.DisplayLayout)
	}
}

func TestLoadConfig_PartialConfigFile(t *testing.T) {
	// Set up a temporary directory for config files
	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	// Create a partial config file
	partialConfig := `
logLevel: debug
anime:
  titleLanguage: native
`

	configPath, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("Failed to get config file path: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	err = os.WriteFile(configPath, []byte(partialConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loading config with partial data
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that specified fields are loaded
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.AnimeConfig.TitleLanguage != "native" {
		t.Errorf("Expected Anime TitleLanguage 'native', got '%s'", cfg.AnimeConfig.TitleLanguage)
	}

	// Check that missing fields have default values
	if cfg.AnimeConfig.DisplayLayout != "list" {
		t.Errorf("Expected default Anime DisplayLayout 'list', got '%s'", cfg.AnimeConfig.DisplayLayout)
	}
}

func TestLoadConfig_WithEnvVar(t *testing.T) {
	// Set up a temporary directory for config files
	tempDir := t.TempDir()
	customConfigPath := filepath.Join(tempDir, "custom_config.yaml")
	os.Setenv("HISAME_CONFIG_FILE", customConfigPath)
	defer os.Unsetenv("HISAME_CONFIG_FILE")

	// Create a custom config file
	customConfig := `
logLevel: warn
`

	err := os.WriteFile(customConfigPath, []byte(customConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom config file: %v", err)
	}

	// Test loading config with env var set
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that specified fields are loaded from custom config
	if cfg.LogLevel != "warn" {
		t.Errorf("Expected LogLevel 'warn', got '%s'", cfg.LogLevel)
	}

	// Check that missing fields have default values
	if cfg.AnimeConfig.TitleLanguage != "english" {
		t.Errorf("Expected default Anime TitleLanguage 'english', got '%s'", cfg.AnimeConfig.TitleLanguage)
	}
	if cfg.AnimeConfig.DisplayLayout != "list" {
		t.Errorf("Expected default Anime DisplayLayout 'list', got '%s'", cfg.AnimeConfig.DisplayLayout)
	}
}

func TestLoadConfig_WithInvalidEnvVarLocation(t *testing.T) {
	// Set up a temporary directory for config files
	tempDir := t.TempDir()
	customConfigPath := filepath.Join(tempDir, "non_existent_config.yaml")
	os.Setenv("HISAME_CONFIG_FILE", customConfigPath)
	defer os.Unsetenv("HISAME_CONFIG_FILE")

	// Ensure the custom config file does not exist
	os.Remove(customConfigPath)

	// Test loading config with env var set to a non-existent file
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config with non-existent env var: %v", err)
	}

	// Check that defaults are used
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default LogLevel 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.AnimeConfig.TitleLanguage != "english" {
		t.Errorf("Expected default Anime TitleLanguage 'english', got '%s'", cfg.AnimeConfig.TitleLanguage)
	}
	if cfg.AnimeConfig.DisplayLayout != "list" {
		t.Errorf("Expected default Anime DisplayLayout 'list', got '%s'", cfg.AnimeConfig.DisplayLayout)
	}
}
