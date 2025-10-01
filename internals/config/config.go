package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName        = "clim_cli"
	ConfigFileName = "clim_cli"
	ConfigFileType = "yaml"
)

// Config represents the application configuration
type Config struct {
	IP      string       `mapstructure:"ip" yaml:"ip"`
	Name    string       `mapstructure:"name" yaml:"name"`
	Power   string       `mapstructure:"power" yaml:"power"`
	Mode    string       `mapstructure:"mode" yaml:"mode"`
	Temp    string       `mapstructure:"temp" yaml:"temp"`
	FanDir  string       `mapstructure:"fan_dir" yaml:"fan_dir"`
	FanRate string       `mapstructure:"fan_rate" yaml:"fan_rate"`
	Search  SearchConfig `mapstructure:"search" yaml:"search"`
}

// SearchConfig represents search-related configuration
type SearchConfig struct {
	Timeout int `mapstructure:"timeout" yaml:"timeout"`
	Workers int `mapstructure:"workers" yaml:"workers"`
}

var configDir string

// InitConfig initializes Viper with the config directory and file
func InitConfig() error {
	// Get config directory
	dir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := ensureConfigDir(dir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configDir = dir

	// Set up Viper
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(dir)

	// Set default values
	setDefaults()

	// Read config file (ignore error if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// GetConfigDir returns the config directory path
func GetConfigDir() string {
	return configDir
}

// getConfigDir returns the appropriate config directory for the platform
func getConfigDir() (string, error) {
	if dir := os.Getenv("CLIM_CLI_CONFIG_DIR"); dir != "" {
		return dir, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, AppName), nil
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// setDefaults sets the default configuration values
func setDefaults() {
	viper.SetDefault("ip", "172.17.2.16")
	viper.SetDefault("name", "")
	viper.SetDefault("power", "1")
	viper.SetDefault("mode", "4")
	viper.SetDefault("temp", "19.0")
	viper.SetDefault("fan_dir", "0")
	viper.SetDefault("fan_rate", "A")
	viper.SetDefault("search.timeout", 5)
	viper.SetDefault("search.workers", 10)
}

// SaveConfig saves the current configuration to file
func SaveConfig() error {
	// Try writing to an existing config file
	if err := viper.WriteConfig(); err != nil {
		// If no existing config, create the default one
		dir := configDir
		if dir == "" {
			var derr error
			dir, derr = getConfigDir()
			if derr != nil {
				return fmt.Errorf("failed to determine config directory: %w", derr)
			}
		}
		if mkerr := ensureConfigDir(dir); mkerr != nil {
			return fmt.Errorf("failed to create config directory: %w", mkerr)
		}
		path := filepath.Join(dir, ConfigFileName+"."+ConfigFileType)
		if werr := viper.WriteConfigAs(path); werr != nil {
			return fmt.Errorf("failed to write config file: %w", werr)
		}
	}
	return nil
}

// SaveConfigAs saves the current configuration to a specific file
func SaveConfigAs(filename string) error {
	return viper.WriteConfigAs(filepath.Join(configDir, filename+"."+ConfigFileType))
}

// LoadConfig loads configuration from a specific file
func LoadConfig(filename string) error {
	viper.SetConfigName(filename)
	return viper.ReadInConfig()
}

// GetDefaultIP returns the default IP address
func GetDefaultIP() string {
	return viper.GetString("ip")
}

// GetDefaultName returns the default device name
func GetDefaultName() string {
	return viper.GetString("name")
}

// GetDefaultPower returns the default power setting
func GetDefaultPower() string {
	return viper.GetString("power")
}

// GetDefaultMode returns the default mode setting
func GetDefaultMode() string {
	return viper.GetString("mode")
}

// GetDefaultTemp returns the default temperature setting
func GetDefaultTemp() string {
	return viper.GetString("temp")
}

// GetDefaultFanDir returns the default fan direction
func GetDefaultFanDir() string {
	return viper.GetString("fan_dir")
}

// GetDefaultFanRate returns the default fan rate
func GetDefaultFanRate() string {
	return viper.GetString("fan_rate")
}

// GetSearchTimeout returns the search timeout
func GetSearchTimeout() int {
	return viper.GetInt("search.timeout")
}

// GetSearchWorkers returns the number of search workers
func GetSearchWorkers() int {
	return viper.GetInt("search.workers")
}

// GetConfig returns the current configuration as a Config struct
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// SetConfig sets the configuration from a Config struct
func SetConfig(cfg *Config) error {
	return viper.MergeConfigMap(map[string]any{
		"ip":       cfg.IP,
		"name":     cfg.Name,
		"power":    cfg.Power,
		"mode":     cfg.Mode,
		"temp":     cfg.Temp,
		"fan_dir":  cfg.FanDir,
		"fan_rate": cfg.FanRate,
		"search": map[string]any{
			"timeout": cfg.Search.Timeout,
			"workers": cfg.Search.Workers,
		},
	})
}

// BindFlags binds cobra command flags to Viper (checks persistent and local)
func BindFlags(cmd *cobra.Command) {
	bind := func(key, name string) {
		if f := cmd.PersistentFlags().Lookup(name); f != nil {
			viper.BindPFlag(key, f)
			return
		}
		if f := cmd.Flags().Lookup(name); f != nil {
			viper.BindPFlag(key, f)
		}
	}
	bind("ip", "ip")
	bind("name", "name")
	bind("power", "power")
	bind("mode", "mode")
	bind("temp", "temp")
	bind("fan_dir", "fan-dir")
	bind("fan_rate", "fan-rate")
}
