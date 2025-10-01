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
	// Initialize config
	if err := config.InitConfig(); err != nil {
		os.Exit(1)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Global flags that bind to Viper config
	rootCmd.PersistentFlags().StringP("ip", "i", config.GetDefaultIP(), "default IP address for climate devices")
	rootCmd.PersistentFlags().StringP("name", "n", config.GetDefaultName(), "default device name")
	rootCmd.PersistentFlags().StringP("power", "p", config.GetDefaultPower(), "default power setting")
	rootCmd.PersistentFlags().StringP("mode", "m", config.GetDefaultMode(), "default mode setting")
	rootCmd.PersistentFlags().StringP("temp", "t", config.GetDefaultTemp(), "default temperature setting")
	rootCmd.PersistentFlags().StringP("fan-dir", "d", config.GetDefaultFanDir(), "default fan direction")
	rootCmd.PersistentFlags().StringP("fan-rate", "r", config.GetDefaultFanRate(), "default fan rate")

	// Bind flags to Viper
	config.BindFlags(rootCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
