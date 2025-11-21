/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/romaingallez/clim_cli/internals/config"
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

	// All flags are now persistent from root command, but can be overridden locally
	setCmd.Flags().StringP("ip", "", "", "IP address (overrides global default)")
	setCmd.Flags().StringP("power", "", "", "power setting (overrides global default)")
	setCmd.Flags().StringP("mode", "", "", "mode setting (overrides global default)")
	setCmd.Flags().StringP("temp", "", "", "temperature setting (overrides global default)")
	setCmd.Flags().StringP("fan-dir", "", "", "fan direction: 0=all wings stopped, 1=vertical, 2=horizontal, 3=both (overrides global default)")
	setCmd.Flags().StringP("fan-rate", "", "", "fan rate (overrides global default)")

	// Bind local flags as well so they override Viper
	config.BindFlags(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
