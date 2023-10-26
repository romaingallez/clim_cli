package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func GetClim(cmd *cobra.Command, args []string) {
	log.Println("clim called")
	// get boolp flag
	toggle, err := cmd.Flags().GetBool("toggle")
	if err != nil {
		log.Println(err)
	}
	if toggle {
		log.Println("toggle is true")
	} else {
		log.Println("toggle is false")
	}
}
