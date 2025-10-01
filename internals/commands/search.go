package commands

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/romaingallez/clim_cli/internals/search"
	"github.com/romaingallez/clim_cli/internals/storage"
	"github.com/romaingallez/clim_cli/internals/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchClim(cmd *cobra.Command, args []string) {
	// Get flags from cobra command
	ifaceName, _ := cmd.Flags().GetString("iface")
	timeout, _ := cmd.Flags().GetInt("timeout")
	workers, _ := cmd.Flags().GetInt("workers")
	tuiMode, _ := cmd.Flags().GetBool("tui")

	fmt.Printf("Searching for climate devices on interface: %s\n", ifaceName)
	fmt.Printf("Timeout: %d seconds, Workers: %d\n", timeout, workers)

	// Use fuzzy search for "*murata*" pattern and save AC manufacturer MACs to config
	devices, err := search.FuzzySearchDevices(ifaceName, timeout, workers, "murata")
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "arp-scan is not installed") {
			fmt.Println("arp-scan not found. Install it first:")
			fmt.Println("  Debian/Ubuntu: sudo apt-get install arp-scan")
			fmt.Println("  macOS (Homebrew): brew install arp-scan")
			return
		}
		if strings.Contains(msg, "interface") && strings.Contains(msg, "not found") {
			fmt.Printf("Network interface '%s' not found. Use --iface to choose a valid interface.\n", ifaceName)
			return
		}
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No climate devices found matching 'murata' pattern")
		return
	}

	// Save devices to storage
	if err := storage.SaveDevices(devices); err != nil {
		log.Printf("Warning: Failed to save devices to storage: %v", err)
	} else {
		fmt.Printf("\nSaved %d device(s) to storage\n", len(devices))
	}

	// If TUI mode is enabled, launch the interactive selector
	if tuiMode {
		fmt.Println("\nLaunching interactive device selector...")
		selectedDevices, err := tui.RunDeviceSelector()
		if err != nil {
			log.Fatalf("Error running TUI: %v", err)
		}

		if len(selectedDevices) > 0 {
			fmt.Printf("\nSelected %d device(s):\n", len(selectedDevices))
			for i, device := range selectedDevices {
				fmt.Printf("%d. %s (%s)\n", i+1, device.Device.Name, device.Device.IP)
				viper.Set("ip", device.Device.IP)
				viper.Set("name", device.Device.Name)
			}

			// Persist the last selected device to the default config file
			if err := config.SaveConfig(); err != nil {
				log.Printf("Warning: Failed to save selected device to config: %v", err)
			} else {
				cfgPath := filepath.Join(config.GetConfigDir(), config.ConfigFileName+"."+config.ConfigFileType)
				fmt.Printf("\nSelected device saved to config: %s\n", cfgPath)
			}
		} else {
			fmt.Println("\nNo devices selected.")
		}
		return
	}

	// Traditional text output
	fmt.Printf("\nFound %d climate device(s) matching 'murata' pattern:\n", len(devices))
	for i, device := range devices {
		fmt.Printf("%d. IP: %s, Status: %s, Name: %s, MAC: %s\n", i+1, device.IP, device.Status, device.Name, device.MAC)

		// Display basic info if available
		if len(device.BasicInfo) > 0 {
			fmt.Printf("   Basic Info: ")
			for key, value := range device.BasicInfo {
				fmt.Printf("%s=%s ", key, value)
			}
			fmt.Println()
		}

		// Display control info if available
		if len(device.ControlInfo) > 0 {
			fmt.Printf("   Control Info: ")
			for key, value := range device.ControlInfo {
				fmt.Printf("%s=%s ", key, value)
			}
			fmt.Println()
		}
	}

	// Display saved AC manufacturer MACs from config
	macs, err := config.GetACManufacturerMACs()
	if err != nil {
		log.Printf("Warning: Could not retrieve saved AC manufacturer MACs: %v", err)
	} else if len(macs) > 0 {
		fmt.Printf("\nSaved AC manufacturer MAC addresses in config:\n")
		for i, mac := range macs {
			fmt.Printf("%d. %s\n", i+1, mac)
		}
	}
}
