package main

import (
	"fmt"
	"os"
	"zkeep/internal/client"

	"github.com/spf13/cobra"
)

var Version = "dev"
var BuildTime = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version and build information for the client and daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Client: %s (built %s)\n", Version, BuildTime)

		c := client.New(client.SocketPath)
		response, err := c.SendCommand("version")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server: %s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
