package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/test/cmd"
)

func main() {
	// Initialize root command
	rootCmd := cmd.NewRootCommand()

	// Execute the root command, which handles subcommands
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
