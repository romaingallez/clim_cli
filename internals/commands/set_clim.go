package commands

import (
	"log"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/spf13/cobra"
)

func SetClim(cmd *cobra.Command, args []string) {
	var err error

	// read flags

	ip, err := cmd.Flags().GetString("ip")
	if err != nil {
		log.Println(err)
	}

	power, err := cmd.Flags().GetString("power")
	if err != nil {
		log.Println(err)
	}

	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		log.Println(err)
	}

	temp, err := cmd.Flags().GetString("temp")
	if err != nil {
		log.Println(err)
	}

	fan_dir, err := cmd.Flags().GetString("fan_dir")
	if err != nil {
		log.Println(err)
	}

	fan_rate, err := cmd.Flags().GetString("fan_rate")
	if err != nil {
		log.Println(err)
	}

	clim := api.Clim{
		IP:      ip,
		Power:   power,
		Mode:    mode,
		Temp:    temp,
		Shum:    "",
		FanDir:  fan_dir,
		FanRate: fan_rate,
	}

	api.Set_Clim(clim)
}
