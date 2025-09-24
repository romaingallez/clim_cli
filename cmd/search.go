/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net"

	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for climate devices on the network",
	Long:  `Search for climate devices on the local network using network discovery`,
	Run:   commands.SearchClim,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	ifaceName, err := getDefaultInterface()
	if err != nil {
		ifaceName = ""
	}
	searchCmd.Flags().StringP("iface", "i", ifaceName, "Network interface to use for device search")
	searchCmd.Flags().IntP("timeout", "t", 5, "timeout in seconds for each device check")
	searchCmd.Flags().IntP("workers", "w", 10, "number of concurrent workers")
}

// getDefaultInterface returns the name of the default network interface
func getDefaultInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
					return iface.Name, nil
				}
			}
		}
	}
	return "", nil
}
