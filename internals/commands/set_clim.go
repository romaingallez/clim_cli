package commands

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

	// Validate IP guidance when empty
	if climConfig.IP == "" {
		fmt.Println("No IP configured. Use --ip, or run 'clim_cli search --tui' or 'clim_cli browse' to select a device.")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.SetClim(ctx, clim); err != nil {
		fmt.Printf("Failed to apply settings to %s: %v\n", clim.IP, err)
		return
	}
	fmt.Printf("Settings applied to %s\n", clim.IP)
}

// getClimConfigFromFlags builds a config from flags, using defaults where not provided
func getClimConfigFromFlags(cmd *cobra.Command) (*config.Config, error) {
	cfg := &config.Config{}

	// Read values from Viper (flags override config automatically via BindPFlag)
	cfg.IP = config.GetDefaultIP()
	cfg.Power = config.GetDefaultPower()
	cfg.Mode = config.GetDefaultMode()
	cfg.Temp = config.GetDefaultTemp()
	cfg.FanDir = config.GetDefaultFanDir()
	cfg.FanRate = config.GetDefaultFanRate()
	cfg.Name = config.GetDefaultName()

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
