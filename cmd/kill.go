package cmd

import (
	"fmt"
	"strconv"

	"github.com/bishalr0y/pman/internal/process"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill <pid>",
	Args:  cobra.ExactArgs(1),
	Short: "Kill a process by PID",
	Long: `Kill a running process using its process ID (PID).

Provide the PID of the process you want to terminate. Use the 'json'
command to find the PID of a process.`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Println("Error: Invalid PID provided")
			return
		}

		if err := process.KillProcessWithPID(int32(pid)); err != nil {
			fmt.Printf("Error: Failed to kill process with PID %d: %v\n", pid, err)
		} else {
			fmt.Printf("Success: Killed process with PID %d\n", pid)
		}
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
