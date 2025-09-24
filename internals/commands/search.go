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

	// Search for climate devices
	devices, err := search.SearchDevices(ifaceName, timeout, workers)
	if err != nil {
		log.Fatalf("Error searching for devices: %v", err)
	}

	if len(devices) == 0 {
		fmt.Println("No climate devices found on the network")
		return
	}

	fmt.Printf("\nFound %d climate device(s):\n", len(devices))
	for i, device := range devices {
		fmt.Printf("%d. IP: %s, Status: %s\n", i+1, device.IP, device.Status)
	}
}
