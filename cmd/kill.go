/*
Copyright © 2026 Bishal Roy <bishalroy895@gmail.com>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/bishalr0y/pman/internal/process"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killForce bool

var killCmd = &cobra.Command{
	Use:          "kill <pid>",
	Args:         cobra.ExactArgs(1),
	Short:        "Kill a process by PID",
	SilenceUsage: true,
	Long: `Kill a running process using its process ID (PID).

By default sends SIGTERM. Use --force for SIGKILL.

Provide the PID of the process you want to terminate. Use the 'json'
command to find the PID of a process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid PID: %w", err)
		}

		if err := process.KillProcessWithPID(int32(pid), killForce); err != nil {
			return fmt.Errorf("failed to kill process with PID %d: %w", pid, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Success: Killed process with PID %d\n", pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().BoolVar(&killForce, "force", false, "Send SIGKILL instead of SIGTERM")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// killCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// killCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
