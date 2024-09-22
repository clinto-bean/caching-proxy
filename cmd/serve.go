/*
Copyright © 2024 Clinton Clark clintobean95@gmail.com
*/
package cmd

import (
	"time"

	"github.com/clinto-bean/caching-proxy/pkg/api"
	"github.com/spf13/cobra"
)

var port, cacheSize, cacheInterval, cacheExpiry int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Starts a new http.Server to handle RESTful calls and cache responses`,
	Run: func(cmd *cobra.Command, args []string) {
		API := api.New(cacheSize, time.Duration(cacheExpiry), time.Duration(cacheInterval))
		API.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	// --port/-p for specifying port to serve requests on
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "port to start the server on")
	// -size/-s for number of items cache can hold
	serveCmd.Flags().IntVarP(&cacheSize, "size", "s", 10, "the number of items the cache can hold")
	// -expiry/-e for time in seconds items will persist in memory
	serveCmd.Flags().IntVarP(&cacheExpiry, "expiry", "e", 600, "the interval of item expiration in seconds")
	// -interval/-i for time in seconds on how often to clean the cache
	serveCmd.Flags().IntVarP(&cacheInterval, "interval", "i", 0, "how often in seconds the cache should check for and clean expired items (default expiry / 10)")
}
