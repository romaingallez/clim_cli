/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set the clim parameters",
	Long:  ``,
	Run:   commands.SetClim,
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("ip", "i", "172.17.2.16", "ip")
	setCmd.Flags().StringP("power", "p", "1", "power")
	setCmd.Flags().StringP("mode", "m", "4", "mode")
	setCmd.Flags().StringP("temp", "t", "19.0", "temp")
	// setCmd.Flags().StringP("shum", "s", "", "shum")
	setCmd.Flags().StringP("fan_dir", "d", "0", "fan_dir")
	setCmd.Flags().StringP("fan_rate", "r", "A", "fan_rate")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
