/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/tui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list stored climate devices",
	Long: `List climate devices stored in local storage.

This command displays a summary of all stored climate devices
including their IP addresses, status, and change history.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.PrintDeviceSummary(); err != nil {
			// Error is already printed by PrintDeviceSummary
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
