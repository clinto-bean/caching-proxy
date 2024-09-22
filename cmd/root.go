/*
Copyright © 2024 Clinton Clark clintobean95@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "A proxy for caching http responses",
	Long:  `This caching proxy is a CLI application which caches http responses. URLs are passed to /fetch on the server and responses are returned to the user as well as cached for a specified amount of time. Type "caching-proxy serve -h" for usage details.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
