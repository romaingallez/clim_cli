package commands

import (
	"fmt"
	"strconv"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/spf13/cobra"
)

func SetClim(cmd *cobra.Command, args []string) {
	// Get configuration with flag overrides
	climConfig, err := getClimConfigFromFlags(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Validate configuration
	if err := validateClimConfig(climConfig); err != nil {
		fmt.Println(err.Error())
		return
	}

	clim := api.Clim{
		IP:      climConfig.IP,
		Power:   climConfig.Power,
		Mode:    climConfig.Mode,
		Temp:    climConfig.Temp,
		Shum:    "",
		FanDir:  climConfig.FanDir,
		FanRate: climConfig.FanRate,
	}

	api.Set_Clim(clim)
}

// getClimConfigFromFlags builds a config from flags, using defaults where not provided
func getClimConfigFromFlags(cmd *cobra.Command) (*config.Config, error) {
	cfg := &config.Config{}

	// Get values from flags, fall back to config defaults
	if ip, _ := cmd.Flags().GetString("ip"); ip != "" {
		cfg.IP = ip
	} else {
		cfg.IP = config.GetDefaultIP()
	}

	if power, _ := cmd.Flags().GetString("power"); power != "" {
		cfg.Power = power
	} else {
		cfg.Power = config.GetDefaultPower()
	}

	if mode, _ := cmd.Flags().GetString("mode"); mode != "" {
		cfg.Mode = mode
	} else {
		cfg.Mode = config.GetDefaultMode()
	}

	if temp, _ := cmd.Flags().GetString("temp"); temp != "" {
		cfg.Temp = temp
	} else {
		cfg.Temp = config.GetDefaultTemp()
	}

	if fanDir, _ := cmd.Flags().GetString("fan-dir"); fanDir != "" {
		cfg.FanDir = fanDir
	} else {
		cfg.FanDir = config.GetDefaultFanDir()
	}

	if fanRate, _ := cmd.Flags().GetString("fan-rate"); fanRate != "" {
		cfg.FanRate = fanRate
	} else {
		cfg.FanRate = config.GetDefaultFanRate()
	}

	if name, _ := cmd.Flags().GetString("name"); name != "" {
		cfg.Name = name
	} else {
		cfg.Name = config.GetDefaultName()
	}

	return cfg, nil
}

// validateClimConfig validates the climate configuration values
func validateClimConfig(cfg *config.Config) error {
	// Validate power
	if cfg.Power != "0" && cfg.Power != "1" {
		return fmt.Errorf("power must be 0 or 1")
	}

	// Validate mode
	if cfg.Mode != "0" && cfg.Mode != "1" && cfg.Mode != "2" && cfg.Mode != "3" && cfg.Mode != "4" {
		return fmt.Errorf("mode must be between 0 and 4")
	}

	// Validate temperature
	tempNum, err := strconv.ParseFloat(cfg.Temp, 64)
	if err != nil {
		return fmt.Errorf("error converting temperature to float: %w", err)
	}
	if tempNum < 16.0 || tempNum > 30.0 {
		return fmt.Errorf("temperature must be between 16.0 and 30.0")
	}

	// Validate fan rate
	if cfg.FanRate != "A" {
		fanNum, err := strconv.Atoi(cfg.FanRate)
		if err != nil {
			return fmt.Errorf("fan_rate must be between 3 and 7 or A")
		}
		if fanNum < 3 || fanNum > 7 {
			return fmt.Errorf("fan_rate must be between 3 and 7 or A")
		}
	}

	return nil
}
