package commands

import (
	"fmt"
	"log"
	"strconv"

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

	// test if power is 0 or 1 (off or on)

	if power != "0" && power != "1" {
		fmt.Println("power must be 0 or 1")
		return
	}

	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		log.Println(err)
	}

	// test if mode is between 0 and 4

	if mode != "0" && mode != "1" && mode != "2" && mode != "3" && mode != "4" {
		fmt.Println("mode must be between 0 and 4")
		return
	}

	temp, err := cmd.Flags().GetString("temp")
	if err != nil {
		log.Println(err)
	}

	// test if temp is between 16.0 and 30.0

	// Convert string to float64
	num, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return
	}

	if num < 16.0 || num > 30.0 {
		fmt.Println("temp must be between 16.0 and 30.0")
		return
	}

	fan_dir, err := cmd.Flags().GetString("fan_dir")
	if err != nil {
		log.Println(err)
	}

	fan_rate, err := cmd.Flags().GetString("fan_rate")
	if err != nil {
		log.Println(err)
	}

	// test if fan_rate is not "A"
	if fan_rate != "A" {
		// convert string to int
		num, err := strconv.Atoi(fan_rate)
		if err != nil {
			// fmt.Println("Error converting string to int:", err)
			fmt.Println("fan_rate must be between 3 and 7 or A")
			return
		}
		// if num is not between 3 and 7
		if num < 3 || num > 7 {
			fmt.Println("fan_rate must be between 3 and 7 or A")
			return
		}

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
