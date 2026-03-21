package cmd

import (
	"fmt"
	"strconv"

	"github.com/bishalr0y/pman/internal/process"

	"github.com/spf13/cobra"
)

// killportCmd represents the killport command
var killportCmd = &cobra.Command{
	Use:   "killport <port>",
	Args:  cobra.ExactArgs(1),
	Short: "Kill a process listening on a specific port.",
	Long:  "Kills the process that is listening on the specified port.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Port argument is required")
			return
		}

		port, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Println("Error: Invalid Port provided")
			return
		}

		if err := process.KillProcessWithPort(int32(port)); err != nil {
			fmt.Printf("Error: Failed to kill process with Port %d: %v\n", port, err)
		} else {
			fmt.Printf("Success: Killed process with Port %d\n", port)
		}
	},
}

func init() {
	rootCmd.AddCommand(killportCmd)
}
