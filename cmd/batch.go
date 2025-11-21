/*
Copyright © 2023 GALLEZ Romain
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/spf13/cobra"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Apply climate settings to multiple devices",
	Long: `Apply climate settings to multiple devices based on group names (grp_name).

Two modes are supported:

1. Simple mode: Apply same settings to all devices in a group
   clim_cli batch --group "coté10" --power 1 --mode 4 --temp 22.0

2. Script mode: Apply settings from a JSON script file (supports multiple groups and per-device overrides)
   clim_cli batch --script ./testdata/batch-example.json`,
	Run: commands.BatchClim,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Script mode flag
	batchCmd.Flags().StringP("script", "s", "", "Path to JSON script file")

	// Simple mode flags
	batchCmd.Flags().StringP("group", "g", "", "Group name (grp_name) to apply settings to")
	batchCmd.Flags().StringP("power", "p", "", "Power setting (0 or 1)")
	batchCmd.Flags().StringP("mode", "m", "", "Mode setting (0=AUTO, 1=HEAT, 2=DRY, 3=FAN, 4=COOL)")
	batchCmd.Flags().StringP("temp", "t", "", "Temperature setting (e.g., 22.0)")
	batchCmd.Flags().StringP("fan-rate", "r", "", "Fan rate (A or 3-7)")
	batchCmd.Flags().StringP("fan-dir", "d", "", "Fan direction (0=all wings stopped, 1=vertical, 2=horizontal, 3=both)")
}

