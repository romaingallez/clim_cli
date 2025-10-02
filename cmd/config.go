/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `View, save, and load configuration settings for the CLI tool.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config: %v\n", err)
			os.Exit(1)
		}

		configFilePath := config.GetConfigFilePath()
		if configFilePath == "" {
			configFilePath = "No config file found (using defaults)"
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("Config File: %s\n", configFilePath)
		fmt.Printf("IP: %s\n", cfg.IP)
		fmt.Printf("Name: %s\n", cfg.Name)
		fmt.Printf("Power: %s\n", cfg.Power)
		fmt.Printf("Mode: %s\n", cfg.Mode)
		fmt.Printf("Temperature: %s\n", cfg.Temp)
		fmt.Printf("Fan Direction: %s\n", cfg.FanDir)
		fmt.Printf("Fan Rate: %s\n", cfg.FanRate)
		fmt.Printf("Search Timeout: %d\n", cfg.Search.Timeout)
		fmt.Printf("Search Workers: %d\n", cfg.Search.Workers)
		fmt.Printf("\nConfig Directory: %s\n", config.GetConfigDir())
	},
}

var configSaveCmd = &cobra.Command{
	Use:   "save [name]",
	Short: "Save current configuration",
	Long:  `Save the current configuration to a file. If name is provided, saves as a named config.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Save to default config
			if err := config.SaveConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Configuration saved to default config file")
		} else {
			// Save as named config
			name := args[0]
			if err := config.SaveConfigAs(name); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config as %s: %v\n", name, err)
				os.Exit(1)
			}
			fmt.Printf("Configuration saved as '%s'\n", name)
		}
	},
}

var configLoadCmd = &cobra.Command{
	Use:   "load <name>",
	Short: "Load a named configuration",
	Long:  `Load a previously saved named configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := config.LoadConfig(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config '%s': %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("Configuration '%s' loaded\n", name)
	},
}

var configDirCmd = &cobra.Command{
	Use:   "dir",
	Short: "Show configuration directory",
	Long:  `Display the path to the configuration directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GetConfigDir())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSaveCmd)
	configCmd.AddCommand(configLoadCmd)
	configCmd.AddCommand(configDirCmd)
}
