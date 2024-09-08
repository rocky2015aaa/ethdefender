package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			// Default action when no subcommands are provided
			fmt.Println("Welcome to testing ETH testnet CLI application. Use --help to see available commands.")
		},
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands
	rootCmd.AddCommand(detectorTest1Cmd)
	rootCmd.AddCommand(detectorTest2Cmd)
	rootCmd.AddCommand(preventerTestCmd)
	rootCmd.AddCommand(upgradeTestCmd)
	rootCmd.AddCommand(unpauseTestCmd)
	rootCmd.AddCommand(checkPaused)

	return rootCmd
}
