/*
Copyright © 2023 GALLEZ Romain
*/
package main

import (
	"log"
	"os"

	"github.com/romaingallez/clim_cli/cmd"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("clim_cli: ")
	log.SetOutput(os.Stderr)

	// add additional code here
}

func main() {

	cmd.Execute()
}
