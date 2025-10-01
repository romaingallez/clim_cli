/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/tui"
	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse stored climate devices",
	Long: `Browse climate devices stored in local storage.

This command launches an interactive TUI to browse and select from
previously discovered climate devices. Devices are sorted by name
and show historical information including change tracking.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunDeviceBrowser(); err != nil {
			// Error is already handled in the TUI
		}
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
