/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get the clim parameters",
	Long:  ``,
	Run:   commands.GetClim,
}

func init() {
	rootCmd.AddCommand(getCmd)

	// IP flag is now persistent from root command, but can be overridden locally
	getCmd.Flags().StringP("ip", "", "", "IP address (overrides global default)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
