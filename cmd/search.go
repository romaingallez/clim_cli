/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net"
	"strings"

	"github.com/romaingallez/clim_cli/internals/commands"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for climate devices on the network",
	Long: `Search for climate devices on the local network using network discovery.

This command discovers climate devices, saves them to local storage with historical
tracking, and optionally launches an interactive TUI for device selection.

Devices are stored with timestamps to track changes over time. Use --tui flag
for interactive selection sorted by device name.`,
	Run: commands.SearchClim,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	ifaceName, err := getDefaultInterface()
	if err != nil {
		ifaceName = ""
	}
	searchCmd.Flags().StringP("iface", "I", ifaceName, "Network interface to use for device search")
	searchCmd.Flags().IntP("timeout", "", 5, "timeout in seconds for each device check")
	searchCmd.Flags().IntP("workers", "w", 10, "number of concurrent workers")
	searchCmd.Flags().Bool("tui", false, "launch interactive TUI for device selection")
}

// getDefaultInterface returns the name of the default network interface
func getDefaultInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Helper function to check if interface is a VPN or virtual interface
	isVPNOrVirtual := func(name string) bool {
		vpnPatterns := []string{"tailscale", "tun", "tap", "wg", "vpn", "docker", "veth", "br-"}
		for _, pattern := range vpnPatterns {
			if strings.HasPrefix(strings.ToLower(name), pattern) {
				return true
			}
		}
		return false
	}

	// Look for physical interfaces (skip VPN and virtual interfaces)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if isVPNOrVirtual(iface.Name) {
				continue
			}
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
