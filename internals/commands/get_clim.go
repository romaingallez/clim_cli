package commands

import (
	"fmt"
	"log"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/romaingallez/clim_cli/internals/config"
	"github.com/spf13/cobra"
)

func GetClim(cmd *cobra.Command, args []string) {
	// Get IP from flag (overrides config default if provided)
	ip, err := cmd.Flags().GetString("ip")
	if err != nil {
		log.Println(err)
		return
	}

	// If no IP provided, use config default
	if ip == "" {
		ip = config.GetDefaultIP()
	}

	basicInfo := api.GetBasicInfo(ip)
	controlInfo := api.GetControlInfo(ip)

	fmt.Printf("Basic info: %+v\nControl Info: %+v\n", basicInfo, controlInfo)
}
