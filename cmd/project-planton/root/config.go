package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BackendURL        string `yaml:"backend-url,omitempty"`
	WebAppContainerID string `yaml:"webapp-container-id,omitempty"`
	WebAppVersion     string `yaml:"webapp-version,omitempty"`
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage openmcf configuration",
	Long:  "Configure openmcf settings like backend URL for API operations",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "set a configuration value",
	Long:  "Set a configuration value. Available keys: backend-url",
	Args:  cobra.ExactArgs(2),
	Run:   configSetHandler,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "get a configuration value",
	Long:  "Get a configuration value. Available keys: backend-url",
	Args:  cobra.ExactArgs(1),
	Run:   configGetHandler,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all configuration values",
	Long:  "List all configuration values",
	Run:   configListHandler,
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".openmcf")
}

func getConfigFile() string {
	return filepath.Join(getConfigDir(), "config.yaml")
}

func loadConfig() (*Config, error) {
	configFile := getConfigFile()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configFile := getConfigFile()
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func configSetHandler(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch key {
	case "backend-url":
		// Basic URL validation
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			fmt.Printf("Error: backend-url must start with http:// or https://\n")
			os.Exit(1)
		}
		config.BackendURL = value
	default:
		fmt.Printf("Error: unknown configuration key '%s'. Available keys: backend-url\n", key)
		os.Exit(1)
	}

	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration %s set to %s\n", key, value)
}

func configGetHandler(cmd *cobra.Command, args []string) {
	key := args[0]

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch key {
	case "backend-url":
		if config.BackendURL == "" {
			fmt.Printf("backend-url is not set\n")
			os.Exit(1)
		}
		fmt.Println(config.BackendURL)
	default:
		fmt.Printf("Error: unknown configuration key '%s'. Available keys: backend-url\n", key)
		os.Exit(1)
	}
}

func configListHandler(cmd *cobra.Command, args []string) {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	hasConfig := false
	if config.BackendURL != "" {
		fmt.Printf("backend-url=%s\n", config.BackendURL)
		hasConfig = true
	}
	if config.WebAppContainerID != "" {
		fmt.Printf("webapp-container-id=%s\n", config.WebAppContainerID)
		hasConfig = true
	}
	if config.WebAppVersion != "" {
		fmt.Printf("webapp-version=%s\n", config.WebAppVersion)
		hasConfig = true
	}

	if !hasConfig {
		fmt.Println("No configuration values set")
	}
}

// GetBackendURL returns the configured backend URL or an error if not set
func GetBackendURL() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	if config.BackendURL == "" {
		return "", fmt.Errorf("backend URL not configured. Run: openmcf config set backend-url <url>")
	}

	return config.BackendURL, nil
}

// LoadConfigPublic loads the configuration (exported for other packages)
func LoadConfigPublic() (*Config, error) {
	return loadConfig()
}

// SaveConfigPublic saves the configuration (exported for other packages)
func SaveConfigPublic(config *Config) error {
	return saveConfig(config)
}
