package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dev_profiler/internal/dto"
	"dev_profiler/internal/utils"
)

const (
	ConfigDirName  = "dev_profiler"
	ConfigFileName = "config.json"
)

// Config holds application configuration
type Config struct {
	GitHub *dto.GitHubConfig `json:"github"`
	OpenAI *dto.OpenAIConfig `json:"openai"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		GitHub: dto.DefaultGitHubConfig(),
		OpenAI: dto.DefaultOpenAIConfig(),
	}
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", ConfigDirName)
	return filepath.Join(configDir, ConfigFileName), nil
}

// GetConfigDir returns the path to the configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", ConfigDirName), nil
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure GitHub config is not nil
	if config.GitHub == nil {
		config.GitHub = dto.DefaultGitHubConfig()
	}

	// Ensure OpenAI config is not nil
	if config.OpenAI == nil {
		config.OpenAI = dto.DefaultOpenAIConfig()
	}

	// Decrypt the GitHub token if it's not empty
	if config.GitHub.Token != "" {
		encryptedData, err := utils.DecodeBase64(config.GitHub.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to decode GitHub token: %w", err)
		}
		decryptedToken, err := utils.DecryptString(encryptedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt GitHub token: %w", err)
		}
		config.GitHub.Token = decryptedToken
	}

	// Decrypt the OpenAI API key if it's not empty
	if config.OpenAI.APIKey != "" {
		encryptedData, err := utils.DecodeBase64(config.OpenAI.APIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode OpenAI API key: %w", err)
		}
		decryptedAPIKey, err := utils.DecryptString(encryptedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt OpenAI API key: %w", err)
		}
		config.OpenAI.APIKey = decryptedAPIKey
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create a copy of config for encryption
	configCopy := *config
	configCopy.GitHub = &dto.GitHubConfig{
		SampledRepoCount:   config.GitHub.SampledRepoCount,
		CommitsPerRepo:     config.GitHub.CommitsPerRepo,
		SampleFileCount:    config.GitHub.SampleFileCount,
		AnalysisYears:      config.GitHub.AnalysisYears,
		IncludePrivateRepo: config.GitHub.IncludePrivateRepo,
		RandomSeed:         config.GitHub.RandomSeed,
		SaveDebugJSON:      config.GitHub.SaveDebugJSON,
	}
	configCopy.OpenAI = &dto.OpenAIConfig{
		Model: config.OpenAI.Model,
	}

	// Encrypt the GitHub token if it's not empty
	if config.GitHub.Token != "" {
		encryptedToken, err := utils.EncryptString(config.GitHub.Token)
		if err != nil {
			return fmt.Errorf("failed to encrypt GitHub token: %w", err)
		}
		configCopy.GitHub.Token = utils.EncodeBase64(encryptedToken)
	} else {
		configCopy.GitHub.Token = ""
	}

	// Encrypt the OpenAI API key if it's not empty
	if config.OpenAI.APIKey != "" {
		encryptedAPIKey, err := utils.EncryptString(config.OpenAI.APIKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt OpenAI API key: %w", err)
		}
		configCopy.OpenAI.APIKey = utils.EncodeBase64(encryptedAPIKey)
	} else {
		configCopy.OpenAI.APIKey = ""
	}

	data, err := json.MarshalIndent(&configCopy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions (owner read/write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
