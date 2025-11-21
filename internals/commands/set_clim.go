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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch current settings from device
	fmt.Printf("Fetching current settings from %s...\n", climConfig.IP)
	currentControlInfo, err := api.FetchControlInfo(ctx, climConfig.IP)
	fetchedCurrent := err == nil
	if err != nil {
		fmt.Printf("Warning: Failed to fetch current settings: %v\n", err)
		fmt.Println("Proceeding with provided values only...")
		currentControlInfo = make(map[string]string)
	}

	// Build current Clim from device response
	currentClim := buildClimFromControlInfo(climConfig.IP, currentControlInfo)

	// Build new Clim: use flag values if provided, otherwise keep current device values or use config defaults
	newClim := buildNewClimFromFlags(cmd, currentClim, climConfig, fetchedCurrent)

	// Validate configuration
	if err := validateClimConfig(&config.Config{
		IP:      newClim.IP,
		Power:   newClim.Power,
		Mode:    newClim.Mode,
		Temp:    newClim.Temp,
		FanDir:  newClim.FanDir,
		FanRate: newClim.FanRate,
	}); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Display current settings
	fmt.Println("\nCurrent settings:")
	displayClimSettings(currentClim)

	// Display new settings and changes
	fmt.Println("\nNew settings:")
	displayClimSettings(newClim)
	displayChanges(currentClim, newClim)

	// Apply new settings
	if err := api.SetClim(ctx, newClim); err != nil {
		fmt.Printf("\nFailed to apply settings to %s: %v\n", newClim.IP, err)
		return
	}
	fmt.Printf("\nSettings applied to %s\n", newClim.IP)
}

// getClimConfigFromFlags builds a config from flags, using defaults where not provided
func getClimConfigFromFlags(cmd *cobra.Command) (*config.Config, error) {
	cfg := &config.Config{}

	// Read flag values directly from command, fallback to config defaults if not provided
	if ip, err := cmd.Flags().GetString("ip"); err == nil && ip != "" {
		cfg.IP = ip
	} else {
		cfg.IP = config.GetDefaultIP()
	}

	if power, err := cmd.Flags().GetString("power"); err == nil && power != "" {
		cfg.Power = power
	} else {
		cfg.Power = config.GetDefaultPower()
	}

	if mode, err := cmd.Flags().GetString("mode"); err == nil && mode != "" {
		cfg.Mode = mode
	} else {
		cfg.Mode = config.GetDefaultMode()
	}

	if temp, err := cmd.Flags().GetString("temp"); err == nil && temp != "" {
		cfg.Temp = temp
	} else {
		cfg.Temp = config.GetDefaultTemp()
	}

	if fanDir, err := cmd.Flags().GetString("fan-dir"); err == nil && fanDir != "" {
		cfg.FanDir = fanDir
	} else {
		cfg.FanDir = config.GetDefaultFanDir()
	}

	if fanRate, err := cmd.Flags().GetString("fan-rate"); err == nil && fanRate != "" {
		cfg.FanRate = fanRate
	} else {
		cfg.FanRate = config.GetDefaultFanRate()
	}

	if name, err := cmd.Flags().GetString("name"); err == nil && name != "" {
		cfg.Name = name
	} else {
		cfg.Name = config.GetDefaultName()
	}

	return cfg, nil
}

// buildClimFromControlInfo builds a Clim struct from the API control info response
func buildClimFromControlInfo(ip string, controlInfo map[string]string) api.Clim {
	clim := api.Clim{
		IP:      ip,
		Power:   controlInfo["pow"],
		Mode:    controlInfo["mode"],
		Temp:    controlInfo["stemp"],
		Shum:    controlInfo["shum"],
		FanRate: controlInfo["f_rate"],
		FanDir:  controlInfo["f_dir"],
	}
	return clim
}

// buildNewClimFromFlags builds a new Clim struct using flag values if provided, otherwise current device values or config defaults
func buildNewClimFromFlags(cmd *cobra.Command, currentClim api.Clim, climConfig *config.Config, fetchedCurrent bool) api.Clim {
	// Start with current device values if we fetched them, otherwise use config defaults
	var power, mode, temp, fanRate, fanDir string
	if fetchedCurrent {
		power = currentClim.Power
		mode = currentClim.Mode
		temp = currentClim.Temp
		fanRate = currentClim.FanRate
		fanDir = currentClim.FanDir
	} else {
		// Use config defaults if we couldn't fetch current settings
		power = climConfig.Power
		mode = climConfig.Mode
		temp = climConfig.Temp
		fanRate = climConfig.FanRate
		fanDir = climConfig.FanDir
	}

	newClim := api.Clim{
		IP:      climConfig.IP,
		Power:   power,
		Mode:    mode,
		Temp:    temp,
		Shum:    currentClim.Shum, // Always use current shum if available, otherwise empty
		FanRate: fanRate,
		FanDir:  fanDir,
	}

	// Override with flag values if provided
	if flagPower, err := cmd.Flags().GetString("power"); err == nil && flagPower != "" {
		newClim.Power = flagPower
	}

	if flagMode, err := cmd.Flags().GetString("mode"); err == nil && flagMode != "" {
		newClim.Mode = flagMode
	}

	if flagTemp, err := cmd.Flags().GetString("temp"); err == nil && flagTemp != "" {
		newClim.Temp = flagTemp
	}

	if flagFanDir, err := cmd.Flags().GetString("fan-dir"); err == nil && flagFanDir != "" {
		newClim.FanDir = flagFanDir
	}

	if flagFanRate, err := cmd.Flags().GetString("fan-rate"); err == nil && flagFanRate != "" {
		newClim.FanRate = flagFanRate
	}

	return newClim
}

// displayClimSettings displays the climate settings in a readable format
func displayClimSettings(clim api.Clim) {
	powerStatus := "OFF"
	if clim.Power == "1" {
		powerStatus = "ON"
	}

	modeNames := map[string]string{
		"0": "AUTO",
		"1": "HEAT",
		"2": "DRY",
		"3": "FAN",
		"4": "COOL",
	}
	modeName := modeNames[clim.Mode]
	if modeName == "" {
		modeName = clim.Mode
	}

	fanDirNames := map[string]string{
		"0": "All wings stopped",
		"1": "Vertical wings motion",
		"2": "Horizontal wings motion",
		"3": "Vertical and horizontal wings motion",
	}
	fanDirName := fanDirNames[clim.FanDir]
	if fanDirName == "" {
		fanDirName = clim.FanDir
	}

	fmt.Printf("  Power:    %s\n", powerStatus)
	fmt.Printf("  Mode:     %s (%s)\n", modeName, clim.Mode)
	fmt.Printf("  Temp:     %s°C\n", clim.Temp)
	fmt.Printf("  Fan Rate: %s\n", clim.FanRate)
	fmt.Printf("  Fan Dir:  %s (%s)\n", fanDirName, clim.FanDir)
}

// displayChanges shows what settings are being changed
func displayChanges(current, new api.Clim) {
	changes := []string{}

	if current.Power != new.Power {
		oldStatus := "OFF"
		if current.Power == "1" {
			oldStatus = "ON"
		}
		newStatus := "OFF"
		if new.Power == "1" {
			newStatus = "ON"
		}
		changes = append(changes, fmt.Sprintf("Power: %s → %s", oldStatus, newStatus))
	}

	if current.Mode != new.Mode {
		changes = append(changes, fmt.Sprintf("Mode: %s → %s", current.Mode, new.Mode))
	}

	if current.Temp != new.Temp {
		changes = append(changes, fmt.Sprintf("Temp: %s°C → %s°C", current.Temp, new.Temp))
	}

	if current.FanRate != new.FanRate {
		changes = append(changes, fmt.Sprintf("Fan Rate: %s → %s", current.FanRate, new.FanRate))
	}

	if current.FanDir != new.FanDir {
		fanDirNames := map[string]string{
			"0": "All wings stopped",
			"1": "Vertical wings motion",
			"2": "Horizontal wings motion",
			"3": "Vertical and horizontal wings motion",
		}
		oldFanDirName := fanDirNames[current.FanDir]
		if oldFanDirName == "" {
			oldFanDirName = current.FanDir
		}
		newFanDirName := fanDirNames[new.FanDir]
		if newFanDirName == "" {
			newFanDirName = new.FanDir
		}
		changes = append(changes, fmt.Sprintf("Fan Dir: %s (%s) → %s (%s)", oldFanDirName, current.FanDir, newFanDirName, new.FanDir))
	}

	if len(changes) > 0 {
		fmt.Println("\nChanges:")
		for _, change := range changes {
			fmt.Printf("  • %s\n", change)
		}
	} else {
		fmt.Println("\nNo changes detected - settings are already as specified.")
	}
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

	// Validate fan direction
	if cfg.FanDir != "0" && cfg.FanDir != "1" && cfg.FanDir != "2" && cfg.FanDir != "3" {
		return fmt.Errorf("fan_dir must be 0 (all wings stopped), 1 (vertical wings motion), 2 (horizontal wings motion), or 3 (vertical and horizontal wings motion)")
	}

	return nil
}
