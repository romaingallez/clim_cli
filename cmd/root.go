/*
Copyright © 2023 GALLEZ Romain
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/romaingallez/clim_cli/internals/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clim_cli",
	Short: "CLI TOOL TO MANAGE CLIM AT WORK",
	Long:  `A golang CLI tool to manage clim at work, using cobra and viper`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag is set
		if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
			v := version.Get()
			fmt.Println(v.Short())
			return
		}

		// If no subcommand and no version flag, show help
		cmd.Help()
	},
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
	rootCmd.PersistentFlags().StringP("fan-dir", "d", "", "default fan direction (0=all wings stopped, 1=vertical, 2=horizontal, 3=both)")
	rootCmd.PersistentFlags().StringP("fan-rate", "r", "", "default fan rate")

	// Version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// Bind flags to Viper
	config.BindFlags(rootCmd)
}
