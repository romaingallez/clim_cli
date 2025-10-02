/*
Copyright © 2023 GALLEZ Romain
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/romaingallez/clim_cli/internals/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print detailed version information including:
- Version number (semantic versioning)
- Git commit hash
- Build date
- Go version used for compilation
- Platform (OS/Architecture)
- Build type (release, development, pre-release)`,
	Run: func(cmd *cobra.Command, args []string) {
		v := version.Get()

		// Check for JSON output flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		shortOutput, _ := cmd.Flags().GetBool("short")

		if jsonOutput {
			fmt.Printf(`{
  "version": "%s",
  "git_commit": "%s",
  "build_date": "%s",
  "go_version": "%s",
  "platform": "%s",
  "build_type": "%s"
}`, v.Version, v.GitCommit, v.BuildDate, v.GoVersion, v.Platform, v.GetBuildType())
			return
		}

		if shortOutput {
			fmt.Println(v.Short())
			return
		}

		// Full output
		fmt.Println(v.String())
		fmt.Printf("Build Type: %s\n", v.GetBuildType())

		// Show warning for development builds
		if v.IsDev() {
			fmt.Fprintf(os.Stderr, "\n⚠️  Warning: This is a development build\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Add flags for different output formats
	versionCmd.Flags().BoolP("json", "j", false, "Output version information in JSON format")
	versionCmd.Flags().BoolP("short", "s", false, "Output only the version number")
}
