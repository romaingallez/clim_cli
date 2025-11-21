package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/romaingallez/clim_cli/internals/storage"
	"github.com/spf13/cobra"
)

// BatchClim handles the batch command to apply settings from a JSON script or simple group flags
func BatchClim(cmd *cobra.Command, args []string) {
	scriptPath, _ := cmd.Flags().GetString("script")
	groupName, _ := cmd.Flags().GetString("group")

	// Determine mode: script mode or simple group mode
	if scriptPath != "" {
		// Script mode
		batchClimFromScript(cmd, scriptPath)
	} else if groupName != "" {
		// Simple group mode
		batchClimFromFlags(cmd, groupName)
	} else {
		fmt.Println("Error: Either --script or --group flag is required")
		fmt.Println("Use --script for JSON script file, or --group for simple group operation")
		return
	}
}

// batchClimFromScript handles batch operations from a JSON script file
func batchClimFromScript(cmd *cobra.Command, scriptPath string) {
	// Read and parse JSON script
	script, err := loadBatchScript(scriptPath)
	if err != nil {
		fmt.Printf("Error loading script: %v\n", err)
		return
	}

	// Load all devices from storage
	allDevices, err := storage.GetDeviceHistories()
	if err != nil {
		fmt.Printf("Error loading devices: %v\n", err)
		return
	}

	if len(allDevices) == 0 {
		fmt.Println("No devices found in storage. Run 'clim_cli search' first.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Process each group configuration
	totalProcessed := 0
	totalSuccess := 0
	totalFailed := 0

	for _, groupConfig := range script.Groups {
		fmt.Printf("\n=== Processing group: %s ===\n", groupConfig.GroupName)

		// Filter devices by grp_name
		matchingDevices := filterDevicesByGroupName(allDevices, groupConfig.GroupName)
		if len(matchingDevices) == 0 {
			fmt.Printf("No devices found with grp_name: %s\n", groupConfig.GroupName)
			continue
		}

		fmt.Printf("Found %d device(s) in group %s\n", len(matchingDevices), groupConfig.GroupName)

		// Process each device
		for _, device := range matchingDevices {
			totalProcessed++

			// Find matching override if any
			override := findDeviceOverride(device, groupConfig.Overrides)

			// Build parameters: start with group defaults, apply override if found
			params := mergeParams(groupConfig.Params, override)

			// Apply settings to device
			success := applySettingsToDevice(ctx, device, params)
			if success {
				totalSuccess++
			} else {
				totalFailed++
			}
		}
	}

	// Summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total devices processed: %d\n", totalProcessed)
	fmt.Printf("Successful: %d\n", totalSuccess)
	fmt.Printf("Failed: %d\n", totalFailed)
}

// batchClimFromFlags handles simple batch operations from command-line flags
func batchClimFromFlags(cmd *cobra.Command, groupName string) {
	// Get parameters from flags
	power, _ := cmd.Flags().GetString("power")
	mode, _ := cmd.Flags().GetString("mode")
	temp, _ := cmd.Flags().GetString("temp")
	fanRate, _ := cmd.Flags().GetString("fan-rate")
	fanDir, _ := cmd.Flags().GetString("fan-dir")

	// Build params from flags
	params := ClimParams{
		Power:   power,
		Mode:    mode,
		Temp:    temp,
		FanRate: fanRate,
		FanDir:  fanDir,
	}

	// Check if at least one parameter is provided
	if power == "" && mode == "" && temp == "" && fanRate == "" && fanDir == "" {
		fmt.Println("Error: At least one parameter (--power, --mode, --temp, --fan-rate, --fan-dir) must be provided")
		return
	}

	// Load all devices from storage
	allDevices, err := storage.GetDeviceHistories()
	if err != nil {
		fmt.Printf("Error loading devices: %v\n", err)
		return
	}

	if len(allDevices) == 0 {
		fmt.Println("No devices found in storage. Run 'clim_cli search' first.")
		return
	}

	// Filter devices by group name
	matchingDevices := filterDevicesByGroupName(allDevices, groupName)
	if len(matchingDevices) == 0 {
		fmt.Printf("No devices found with grp_name: %s\n", groupName)
		return
	}

	fmt.Printf("Found %d device(s) in group: %s\n", len(matchingDevices), groupName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Process each device
	totalProcessed := 0
	totalSuccess := 0
	totalFailed := 0

	for _, device := range matchingDevices {
		totalProcessed++
		success := applySettingsToDevice(ctx, device, params)
		if success {
			totalSuccess++
		} else {
			totalFailed++
		}
	}

	// Summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total devices processed: %d\n", totalProcessed)
	fmt.Printf("Successful: %d\n", totalSuccess)
	fmt.Printf("Failed: %d\n", totalFailed)
}

// loadBatchScript loads and parses the JSON script file
func loadBatchScript(path string) (*BatchScript, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read script file: %w", err)
	}

	var script BatchScript
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(script.Groups) == 0 {
		return nil, fmt.Errorf("script contains no groups")
	}

	return &script, nil
}

// filterDevicesByGroupName filters devices by grp_name from basic_info
func filterDevicesByGroupName(devices []*storage.DeviceHistory, groupName string) []*storage.DeviceHistory {
	var filtered []*storage.DeviceHistory
	for _, device := range devices {
		if grpName, ok := device.Device.BasicInfo["grp_name"]; ok && grpName == groupName {
			filtered = append(filtered, device)
		}
	}
	return filtered
}

// findDeviceOverride finds a matching override for a device by name or IP
func findDeviceOverride(device *storage.DeviceHistory, overrides []DeviceOverride) *ClimParams {
	for _, override := range overrides {
		// Match by name if specified
		if override.Name != "" && device.Device.Name == override.Name {
			return &override.Params
		}
		// Match by IP if specified
		if override.IP != "" && device.Device.IP == override.IP {
			return &override.Params
		}
		// Match if both name and IP are specified and both match
		if override.Name != "" && override.IP != "" {
			if device.Device.Name == override.Name && device.Device.IP == override.IP {
				return &override.Params
			}
		}
	}
	return nil
}

// mergeParams merges default params with override params
// Override takes precedence, but empty strings in override mean "keep default"
func mergeParams(defaults ClimParams, override *ClimParams) ClimParams {
	result := defaults
	if override == nil {
		return result
	}

	// Override only non-empty values
	if override.Power != "" {
		result.Power = override.Power
	}
	if override.Mode != "" {
		result.Mode = override.Mode
	}
	if override.Temp != "" {
		result.Temp = override.Temp
	}
	if override.FanRate != "" {
		result.FanRate = override.FanRate
	}
	if override.FanDir != "" {
		result.FanDir = override.FanDir
	}

	return result
}

// applySettingsToDevice applies settings to a single device
// Returns true on success, false on failure
func applySettingsToDevice(ctx context.Context, device *storage.DeviceHistory, params ClimParams) bool {
	deviceIP := device.Device.IP
	deviceName := device.Device.Name

	fmt.Printf("\n  Device: %s (%s)\n", deviceName, deviceIP)

	// Fetch current settings
	currentControlInfo, err := api.FetchControlInfo(ctx, deviceIP)
	if err != nil {
		fmt.Printf("    Error: Failed to fetch current settings: %v\n", err)
		return false
	}

	// Build current Clim from device response
	currentClim := buildClimFromControlInfoForBatch(deviceIP, currentControlInfo)

	// Build new Clim: merge current settings with script params
	// Empty string in params means "keep current value"
	newClim := api.Clim{
		IP:      deviceIP,
		Power:   getValueOrDefault(params.Power, currentClim.Power),
		Mode:    getValueOrDefault(params.Mode, currentClim.Mode),
		Temp:    getValueOrDefault(params.Temp, currentClim.Temp),
		Shum:    currentClim.Shum, // Always keep current shum
		FanRate: getValueOrDefault(params.FanRate, currentClim.FanRate),
		FanDir:  getValueOrDefault(params.FanDir, currentClim.FanDir),
	}

	// Display what will be changed
	displayChangesBatch(currentClim, newClim)

	// Apply new settings
	if err := api.SetClim(ctx, newClim); err != nil {
		fmt.Printf("    Error: Failed to apply settings: %v\n", err)
		return false
	}

	fmt.Printf("    ✓ Settings applied successfully\n")
	return true
}

// getValueOrDefault returns the value if not empty, otherwise returns the default
func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// buildClimFromControlInfoForBatch builds a Clim struct from the API control info response
func buildClimFromControlInfoForBatch(ip string, controlInfo map[string]string) api.Clim {
	return api.Clim{
		IP:      ip,
		Power:   controlInfo["pow"],
		Mode:    controlInfo["mode"],
		Temp:    controlInfo["stemp"],
		Shum:    controlInfo["shum"],
		FanRate: controlInfo["f_rate"],
		FanDir:  controlInfo["f_dir"],
	}
}

// displayChangesBatch shows what settings are being changed
// (Similar to displayChanges in set_clim.go but adapted for batch output)
func displayChangesBatch(current, new api.Clim) {
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
		fmt.Printf("    Changes: ")
		for i, change := range changes {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", change)
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("    No changes needed\n")
	}
}

