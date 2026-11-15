/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/bishalr0y/pman/internal/process"

	"github.com/spf13/cobra"
)

var killportForce bool

// killportCmd represents the killport command
var killportCmd = &cobra.Command{
	Use:          "killport <port>",
	Args:         cobra.ExactArgs(1),
	Short:        "Kill a process listening on a specific port.",
	SilenceUsage: true,
	Long: `Kills the process that is listening on the specified port.

By default sends SIGTERM. Use --force for SIGKILL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid port: %w", err)
		}

		if err := process.KillProcessWithPort(int32(port), killportForce); err != nil {
			return fmt.Errorf("failed to kill process on port %d: %w", port, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Success: Killed process with Port %d\n", port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(killportCmd)
	killportCmd.Flags().BoolVar(&killportForce, "force", false, "Send SIGKILL instead of SIGTERM")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// killportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// killportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
