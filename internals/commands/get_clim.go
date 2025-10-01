package commands

import (
	"context"
	"fmt"
	"log"
	"time"

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

	if ip == "" {
		fmt.Println("No IP configured. Use --ip, or run 'clim_cli search --tui' or 'clim_cli browse' to select a device.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	basicInfo, berr := api.FetchBasicInfo(ctx, ip)
	controlInfo, cerr := api.FetchControlInfo(ctx, ip)
	if berr != nil {
		fmt.Printf("Failed to fetch basic_info from %s: %v\n", ip, berr)
	}
	if cerr != nil {
		fmt.Printf("Failed to fetch control_info from %s: %v\n", ip, cerr)
	}
	if berr == nil {
		fmt.Printf("Basic info: %+v\n", basicInfo)
	}
	if cerr == nil {
		fmt.Printf("Control Info: %+v\n", controlInfo)
	}
}
