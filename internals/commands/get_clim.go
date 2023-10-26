package commands

import (
	"fmt"
	"log"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/spf13/cobra"
)

func GetClim(cmd *cobra.Command, args []string) {

	ip, err := cmd.Flags().GetString("ip")
	if err != nil {
		log.Println(err)
	}

	basicInfo := api.GetBasicInfo(ip)

	controlInfo := api.GetControlInfo(ip)

	fmt.Printf("Basic info: %+v\nControl Info: %+v\n", basicInfo, controlInfo)
}
