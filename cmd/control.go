package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/romaingallez/clim_cli/internals/storage"
	"github.com/romaingallez/clim_cli/internals/tui"
	"github.com/spf13/cobra"
)

var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Interactive control TUI for power/mode/temp/fan",
	Run: func(cmd *cobra.Command, args []string) {
		useSelector, _ := cmd.Flags().GetBool("tui-select")
		ipsArg, _ := cmd.Flags().GetString("ips")

		var selected []*storage.DeviceHistory
		var err error

		if useSelector {
			selected, err = tui.RunDeviceSelector()
			if err != nil {
				log.Fatalf("Error running device selector: %v", err)
			}
			if len(selected) == 0 {
				fmt.Println("No devices selected.")
				return
			}
		} else if ipsArg != "" {
			ips := strings.Split(ipsArg, ",")
			// load histories and match by IP
			all, err := storage.GetDeviceHistories()
			if err != nil {
				log.Fatalf("Failed to load devices: %v", err)
			}
			for _, ip := range ips {
				ip = strings.TrimSpace(ip)
				for _, dh := range all {
					if dh.Device.IP == ip {
						selected = append(selected, dh)
						break
					}
				}
			}
			if len(selected) == 0 {
				fmt.Println("Provided IPs not found in storage.")
				return
			}
		} else {
			fmt.Println("No selection provided. Use --tui-select or --ips.")
			return
		}

		if err := tui.RunControlScreen(selected); err != nil {
			log.Fatalf("Error running control TUI: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(controlCmd)
	controlCmd.Flags().Bool("tui-select", false, "Open device selector before control screen")
	controlCmd.Flags().String("ips", "", "Comma-separated IPs to control (must exist in storage)")
}
