package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "CLI application",
	Long:  "Command Line Interface for clean architect application",
}

func init() {
	rootCmd.AddCommand(
		routingCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
