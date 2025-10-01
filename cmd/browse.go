/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/romaingallez/clim_cli/internals/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse and select stored climate devices",
	Long: `Browse and select climate devices stored in local storage.

This command launches an interactive TUI to browse and select from
previously discovered climate devices. The selected device will be
saved to the current configuration. Devices are sorted by name
and show historical information including change tracking.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedDevices, err := tui.RunDeviceSelector()
		if err != nil {
			log.Fatalf("Error running device selector: %v", err)
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
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
