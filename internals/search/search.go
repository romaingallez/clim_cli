package search

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/spf13/viper"
)

// Device represents a found climate device
type Device struct {
	IP          string
	Status      string
	Name        string
	Model       string
	MAC         string
	BasicInfo   map[string]string
	ControlInfo map[string]string
}

// SearchDevices searches for climate devices in the specified IP range using arp-scan
func SearchDevices(ifaceString string, timeout int, workers int) ([]Device, error) {
	log.Printf("Searching on interface: %s, timeout: %d, workers: %d", ifaceString, timeout, workers)

	// Validate interface exists
	iface, err := net.InterfaceByName(ifaceString)
	if err != nil {
		return nil, fmt.Errorf("interface %s not found: %v", ifaceString, err)
	}

	// Get the network address for the interface
	networkAddr, err := getNetworkAddress(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to get network address for interface %s: %v", ifaceString, err)
	}

	// Execute arp-scan command
	devices, err := executeArpScan(ifaceString, networkAddr, timeout)
	if err != nil {
		return nil, fmt.Errorf("arp-scan failed: %v", err)
	}

	log.Printf("Found %d devices", len(devices))
	return devices, nil
}

// getNetworkAddress returns the network address for the given interface
func getNetworkAddress(iface *net.Interface) (string, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				// Skip localhost
				if ip4[0] == 127 {
					continue
				}
				// Return the network in CIDR notation
				return ipNet.String(), nil
			}
		}
	}

	return "", errors.New("no valid IPv4 address found on interface")
}

// executeArpScan runs the arp-scan command and parses its output
func executeArpScan(iface string, networkAddr string, timeout int) ([]Device, error) {
	// Build the arp-scan command
	cmd := exec.Command("arp-scan", "--localnet", "-x", "--format=${ip};${mac};${vendor}")

	// switch {
	// case iface != "" && networkAddr != "":
	// 	cmd.Args = append(cmd.Args, "-I", iface, "--interface-range", networkAddr)
	// case iface != "":
	// 	cmd.Args = append(cmd.Args, "-I", iface)
	// case networkAddr != "":
	// 	cmd.Args = append(cmd.Args, "--interface-range", networkAddr)
	// }

	// Set timeout

	cmd.Args = append(cmd.Args, "--timeout", fmt.Sprintf("%d", timeout))

	log.Printf("Executing command: %s", strings.Join(cmd.Args, " "))

	// Execute the command
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute arp-scan: %v", err)
	}

	// Parse the output
	devices, err := parseArpScanOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse arp-scan output: %v", err)
	}

	return devices, nil
}

// parseArpScanOutput parses the output from arp-scan command
// Expected format: IP;MAC;Vendor
func parseArpScanOutput(output string) ([]Device, error) {
	var devices []Device

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and header lines
		if line == "" || strings.Contains(line, "Starting arp-scan") ||
			strings.Contains(line, "Interface:") || strings.Contains(line, "Starting:") ||
			strings.Contains(line, "Ending arp-scan") || strings.Contains(line, "packets") {
			continue
		}

		// Parse the line: IP;MAC;Vendor
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue // Skip malformed lines
		}

		ip := strings.TrimSpace(parts[0])
		mac := strings.TrimSpace(parts[1])
		vendor := ""
		if len(parts) > 2 {
			vendor = strings.TrimSpace(parts[2])
		}

		// Validate IP address
		if net.ParseIP(ip) == nil {
			continue
		}

		device := Device{
			IP:     ip,
			MAC:    mac,
			Name:   vendor,
			Status: "online",
		}

		devices = append(devices, device)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// FuzzySearchDevices searches for devices matching a fuzzy pattern and saves AC manufacturer MACs to config
func FuzzySearchDevices(ifaceString string, timeout int, workers int, pattern string) ([]Device, error) {
	// First, get all devices
	devices, err := SearchDevices(ifaceString, timeout, workers)
	if err != nil {
		return nil, err
	}

	// Filter devices using fuzzy search
	filteredDevices := fuzzyFilter(devices, pattern)

	// Get clim info for each filtered device
	for i := range filteredDevices {
		device := &filteredDevices[i]
		log.Printf("Getting clim info for device %s", device.IP)

		// Get basic info
		basicInfo := api.GetBasicInfo(device.IP)
		if basicInfo != nil {
			device.BasicInfo = basicInfo
			// Update device name from basic info if available
			if name, exists := basicInfo["name"]; exists && name != "" {
				device.Name = name
			}
		}

		// Get control info
		controlInfo := api.GetControlInfo(device.IP)
		if controlInfo != nil {
			device.ControlInfo = controlInfo
		}
	}

	// Save AC manufacturer MACs to config
	if err := saveACManufacturerMACs(filteredDevices); err != nil {
		log.Printf("Warning: Failed to save AC manufacturer MACs to config: %v", err)
	}

	return filteredDevices, nil
}

// fuzzyFilter filters devices based on fuzzy pattern matching
func fuzzyFilter(devices []Device, pattern string) []Device {
	var filtered []Device

	// Convert pattern to lowercase for case-insensitive matching
	pattern = strings.ToLower(pattern)

	for _, device := range devices {
		// Check if device name contains the pattern (case-insensitive)
		if strings.Contains(strings.ToLower(device.Name), pattern) {
			filtered = append(filtered, device)
		}
	}

	return filtered
}

// initializeViperConfig initializes Viper with the correct config path
func initializeViperConfig() (string, error) {
	// Initialize Viper
	viper.SetConfigName("clim_cli")
	viper.SetConfigType("yaml")

	// Set config path to user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "clim_cli")
	viper.AddConfigPath(configDir)

	return configDir, nil
}

// saveACManufacturerMACs saves AC manufacturer MAC addresses to Viper config
func saveACManufacturerMACs(devices []Device) error {
	configDir, err := initializeViperConfig()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Set the config file path explicitly
	configFile := filepath.Join(configDir, "clim_cli.yaml")
	viper.SetConfigFile(configFile)

	// Try to read existing config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %v", err)
		}
		// Config file doesn't exist, that's okay - we'll create it
		log.Printf("Config file doesn't exist, will create new one at: %s", configFile)
	}

	// Extract AC manufacturer MACs (devices with "murata" in name)
	var acMACs []string
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Name), "murata") {
			acMACs = append(acMACs, device.MAC)
		}
	}

	// Save AC manufacturer MACs to config
	viper.Set("ac_manufacturer_macs", acMACs)

	// Write config file (this will create the file if it doesn't exist)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Printf("Saved %d AC manufacturer MAC addresses to config", len(acMACs))
	return nil
}

// GetACManufacturerMACs retrieves AC manufacturer MAC addresses from config
func GetACManufacturerMACs() ([]string, error) {
	_, err := initializeViperConfig()
	if err != nil {
		return nil, err
	}

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return []string{}, nil // No config file, return empty slice
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Get AC manufacturer MACs
	macs := viper.GetStringSlice("ac_manufacturer_macs")
	return macs, nil
}
