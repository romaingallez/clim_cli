/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clim_cli",
	Short: "CLI TOOL TO MANAGE CLIM AT WORK",
	Long:  `A golang CLI tool to manage clim at work, using cobra and viper`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize config before defining flags so defaults come from file
	if err := config.InitConfig(); err != nil {
		os.Exit(1)
	}

	// Global flags that bind to Viper config
	// Define flags with empty defaults; Viper provides values
	rootCmd.PersistentFlags().StringP("ip", "i", "", "default IP address for climate devices")
	rootCmd.PersistentFlags().StringP("name", "n", "", "default device name")
	rootCmd.PersistentFlags().StringP("power", "p", "", "default power setting")
	rootCmd.PersistentFlags().StringP("mode", "m", "", "default mode setting")
	rootCmd.PersistentFlags().StringP("temp", "t", "", "default temperature setting")
	rootCmd.PersistentFlags().StringP("fan-dir", "d", "", "default fan direction")
	rootCmd.PersistentFlags().StringP("fan-rate", "r", "", "default fan rate")

	// Bind flags to Viper
	config.BindFlags(rootCmd)
}
