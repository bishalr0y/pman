package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// gets update by github action
const Version = "0.8.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the pman version",
	Long:  `Print the version number and build information of pman.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pman version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
