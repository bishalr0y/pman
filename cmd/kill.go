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

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill a process by PID",
	Long: `Kill a running process using its process ID (PID).

Provide the PID of the process you want to terminate. Use the 'json'
command to find the PID of a process.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: PID argument is required")
			return
		}

		pid, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Println("Error: Invalid PID provided")
			return
		}

		if err := process.KillProcess(int32(pid)); err != nil {
			fmt.Printf("Error: Failed to kill process with PID %d: %v\n", pid, err)
		} else {
			fmt.Printf("Success: Killed process with PID %d\n", pid)
		}
	},
}

func init() {
	rootCmd.AddCommand(killCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// killCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// killCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
