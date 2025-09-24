package commands

import (
	"fmt"
	"log"

	"github.com/romaingallez/clim_cli/internals/search"
	"github.com/spf13/cobra"
)

func SearchClim(cmd *cobra.Command, args []string) {
	// Get flags from cobra command
	ifaceName, _ := cmd.Flags().GetString("iface")
	timeout, _ := cmd.Flags().GetInt("timeout")
	workers, _ := cmd.Flags().GetInt("workers")

	fmt.Printf("Searching for climate devices on interface: %s\n", ifaceName)
	fmt.Printf("Timeout: %d seconds, Workers: %d\n", timeout, workers)

	// Use fuzzy search for "*murata*" pattern and save AC manufacturer MACs to config
	devices, err := search.FuzzySearchDevices(ifaceName, timeout, workers, "murata")
	if err != nil {
		log.Fatalf("Error searching for devices: %v", err)
	}

	if len(devices) == 0 {
		fmt.Println("No climate devices found matching 'murata' pattern")
		return
	}

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
	macs, err := search.GetACManufacturerMACs()
	if err != nil {
		log.Printf("Warning: Could not retrieve saved AC manufacturer MACs: %v", err)
	} else if len(macs) > 0 {
		fmt.Printf("\nSaved AC manufacturer MAC addresses in config:\n")
		for i, mac := range macs {
			fmt.Printf("%d. %s\n", i+1, mac)
		}
	}
}
