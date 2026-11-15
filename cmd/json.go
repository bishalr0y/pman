/*
Copyright © 2026 Bishal Roy <bishalroy895@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/bishalr0y/pman/internal/process"
	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:          "json",
	Short:        "List running processes in JSON format",
	SilenceUsage: true,
	Long: `Display all running processes as a JSON array.

This command retrieves information about all currently running processes
on the system and outputs them in JSON format for easy processing
by other tools or scripts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		processes, err := process.ListProcesses()
		if err != nil {
			return fmt.Errorf("failed to get processes: %w", err)
		}

		jsonProcesses, err := json.MarshalIndent(processes, "", " ")
		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(jsonProcesses))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
